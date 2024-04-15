//go:build integration
// +build integration

package default_controller

import (
	"fmt"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestUpdateClassroom(t *testing.T) {
	testDb := tests.NewTestDB(t)

	user := database.User{ID: 1}
	testDb.InsertUser(&user)

	classroom := database.Classroom{
		ID:          tests.UUID("00000000-0000-0000-0000-000000000001"),
		Name:        "Test",
		Description: "test",
		OwnerID:     user.ID,
		GroupID:     2,
	}
	testDb.InsertClassroom(&classroom)

	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		fiberContext.SetGitlabRepository(c, gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(createdUser.ID)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("UpdateClassroom", func(t *testing.T) {
		app.Post("/api/classrooms", handler.UpdateClassroom)

		requestBody := UpdateClassroomRequest{
			Name:        classroom.Name + "_New",
			Description: classroom.Name + "_new",
		}

		gitlabRepo.
			EXPECT().
			ChangeGroupName(
				classroom.GroupID,
				requestBody.Name,
			).
			Return(
				&model.Group{
					Name: requestBody.Name,
				},
				nil,
			).
			Times(1)

		gitlabRepo.
			EXPECT().
			ChangeGroupDescription(
				classroom.GroupID,
				requestBody.Description,
			).
			Return(
				&model.Group{
					Description: requestBody.Description,
				},
				nil,
			).
			Times(1)

		req := newPutJsonRequest("/api/classrooms/owned/:classroomId", requestBody)
		resp, err := app.Test(req)
		assert.Equal(t, fiber.StatusAccepted, resp.StatusCode)
		assert.NoError(t, err)

		classRoom, err := query.Classroom.Where(query.Classroom.OwnerID.Eq(user.ID)).First()
		assert.NoError(t, err)
		assert.Equal(t, requestBody.Name, classRoom.Name)
		assert.Equal(t, requestBody.Description, classRoom.Description)

		assert.Equal(t, fmt.Sprintf("/api/v1/classrooms/owned/%s", classRoom.ID.String()), resp.Header.Get("Location"))
	})
}
