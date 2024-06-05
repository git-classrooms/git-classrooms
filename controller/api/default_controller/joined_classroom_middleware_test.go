package default_controller

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestJoinedClassroomMiddleware(t *testing.T) {
	testDB := db_tests.NewTestDB(t)

	user := factory.User()
	testDB.InsertUser(&user)

	classroom := factory.Classroom()
	testDB.InsertClassroom(&classroom)

	userClassroom := factory.UserClassroom(user.ID, classroom.ID)
	testDB.InsertUserClassroom(&userClassroom)

	userClassroom.Classroom = classroom
	userClassroom.User = user

	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api/v1/", func(c *fiber.Ctx) error {
		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("JoinedClassroomMiddleware", func(t *testing.T) {
		app.Use("/api/v1/classrooms/joined/:classroomId", handler.JoinedClassroomMiddleware)
		app.Get("/api/v1/classrooms/joined/:classroomId", func(c *fiber.Ctx) error {
			ctx := context.Get(c)
			joinedClassroom := ctx.GetJoinedClassroom()

			assert.Equal(t, classroom.Name, joinedClassroom.Classroom.Name)
			assert.Equal(t, classroom.Description, joinedClassroom.Classroom.Description)
			assert.Equal(t, classroom.GroupID, joinedClassroom.Classroom.GroupID)
			assert.Equal(t, classroom.GroupAccessTokenID, joinedClassroom.Classroom.GroupAccessTokenID)
			assert.Equal(t, classroom.GroupAccessToken, joinedClassroom.Classroom.GroupAccessToken)

			return c.JSON(nil)
		})

		route := fmt.Sprintf("/api/v1/classrooms/joined/%s", classroom.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

	})
}
