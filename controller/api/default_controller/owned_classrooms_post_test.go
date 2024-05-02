//go:build integration
// +build integration

package default_controller

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
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

func TestCreateClassroom(t *testing.T) {
	testDb := tests.NewTestDB(t)

	user := database.User{ID: 1}
	testDb.InsertUser(&user)

	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		fiberContext.SetGitlabRepository(c, gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(user.ID)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("CreateClassroom", func(t *testing.T) {
		app.Post("/api/classrooms", handler.CreateClassroom)

		requestBody := CreateClassroomRequest{
			Name:                    "Test",
			MemberEmails:            []string{},
			Description:             "test",
			StudentsViewAllProjects: true,
		}

		gitlabRepo.
			EXPECT().
			CreateGroup(
				requestBody.Name,
				model.Private,
				requestBody.Description,
			).
			Return(
				&model.Group{ID: 1},
				nil,
			).
			Times(1)

		gitlabRepo.
			EXPECT().
			CreateGroupAccessToken(
				1,
				"Gitlab Classrooms",
				model.OwnerPermissions,
				mock.AnythingOfType("time.Time"),
				"api",
			).
			Return(
				&model.GroupAccessToken{ID: 20, Token: "token"},
				nil,
			).
			Times(1)

		req := newPostJsonRequest("/api/classrooms", requestBody)
		resp, err := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		assert.NoError(t, err)

		classRoom, err := query.Classroom.Where(query.Classroom.OwnerID.Eq(user.ID)).First()
		assert.NoError(t, err)
		assert.Equal(t, "Test", classRoom.Name)
		assert.Equal(t, "test", classRoom.Description)
		assert.Equal(t, 1, classRoom.GroupID)
		assert.Equal(t, 20, classRoom.GroupAccessTokenID)
		assert.Equal(t, "token", classRoom.GroupAccessToken)
		assert.Equal(t, true, classRoom.StudentsViewAllProjects)

		assert.Equal(t, fmt.Sprintf("/api/v1/classrooms/%s", classRoom.ID.String()), resp.Header.Get("Location"))
	})
}
