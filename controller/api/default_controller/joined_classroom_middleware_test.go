package default_controller

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestJoinedClassroomMiddleware(t *testing.T) {
	testDB := db_tests.NewTestDB(t)

	user := factory.User()
	testDB.InsertUser(user)

	classroom := factory.Classroom()
	testDB.InsertClassroom(classroom)

	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetJoinedClassroom(&database.UserClassrooms{ClassroomID: classroom.ID})

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("JoinedClassroomMiddleware", func(t *testing.T) {
		app.Use("/api/classrooms/joined/:classroomId", handler.JoinedClassroomMiddleware)
		route := fmt.Sprintf("/api/classrooms/joined/%s", classroom.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		joinedClassroom, err := query.UserClassrooms.WithContext(context.Background()).
			Preload(query.UserClassrooms.Classroom).
			Where(query.UserClassrooms.UserID.Eq(1)).
			First()

		assert.NoError(t, err)
		assert.Equal(t, classroom.Name, joinedClassroom.Classroom.Name)
		assert.Equal(t, classroom.Description, joinedClassroom.Classroom.Description)
		assert.Equal(t, classroom.GroupID, joinedClassroom.Classroom.GroupID)
		assert.Equal(t, classroom.GroupAccessTokenID, joinedClassroom.Classroom.GroupAccessTokenID)
		assert.Equal(t, classroom.GroupAccessToken, joinedClassroom.Classroom.GroupAccessToken)
	})
}
