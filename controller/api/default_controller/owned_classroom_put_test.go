package default_controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	contextWrapper "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPutOwnedClassroom(t *testing.T) {
	testDb := db_tests.NewTestDB(t)

	user := database.User{ID: 1}
	testDb.InsertUser(&user)

	classroom := database.Classroom{
		ID:          uuid.UUID{},
		Name:        "Test",
		Description: "test",
		OwnerID:     user.ID,
		GroupID:     228,
	}
	testDb.InsertClassroom(&classroom)

	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := contextWrapper.Get(c)
		ctx.SetOwnedClassroom(&classroom)
		ctx.SetGitlabRepository(gitlabRepo)

		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(user.ID)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)
	app.Put("/api/classrooms/owned/:classroomId", handler.PutOwnedClassroom)

	targetRoute := fmt.Sprintf("/api/classrooms/owned/%s", classroom.ID.String())

	t.Run("updates classroom", func(t *testing.T) {
		requestBody := updateClassroomRequest{
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

		req := newPutJsonRequest(targetRoute, requestBody)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusAccepted, resp.StatusCode)
		assert.NoError(t, err)

		classRoom, err := query.Classroom.WithContext(context.Background()).Where(query.Classroom.OwnerID.Eq(user.ID)).First()

		assert.NoError(t, err)
		assert.Equal(t, requestBody.Name, classRoom.Name)
		assert.Equal(t, requestBody.Description, classRoom.Description)
	})

	t.Run("request body is emtpy", func(t *testing.T) {
		requestBody := updateClassroomRequest{}

		req := newPutJsonRequest(targetRoute, requestBody)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		bodyString := string(bodyBytes)
		assert.Equal(t, "request requires name and description", bodyString)
	})

	t.Run("gitlab error on change group name", func(t *testing.T) {
		requestBody := updateClassroomRequest{
			Name:        classroom.Name + "_New",
			Description: classroom.Name + "_new",
		}

		errMsg := "error"
		gitlabRepo.
			EXPECT().
			ChangeGroupName(
				classroom.GroupID,
				requestBody.Name,
			).
			Return(
				nil,
				fmt.Errorf(errMsg),
			).
			Times(1)

		req := newPutJsonRequest(targetRoute, requestBody)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("gitlab error on change group description", func(t *testing.T) {
		requestBody := updateClassroomRequest{
			Name:        classroom.Name + "_New",
			Description: classroom.Name + "_new",
		}

		errMsg := "error"
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
				nil,
				fmt.Errorf(errMsg),
			).
			Times(1)

		req := newPutJsonRequest(targetRoute, requestBody)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}
