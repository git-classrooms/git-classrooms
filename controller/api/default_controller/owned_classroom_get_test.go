package default_controller

import (
	"context"
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetOwnedClassroom(t *testing.T) {
	// --------------- DB SETUP -----------------
	t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	pq, err := tests.StartPostgres()

	user := factory.User()
	testDB.InsertUser(&user)

	classroom := factory.Classroom()
	testDB.InsertClassroom(&classroom)

	session.InitSessionStore(nil, &url.URL{Scheme: "http"})
	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetOwnedClassroom(&classroom)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("GetOwnedClassroom", func(t *testing.T) {
		app.Get("/api/classrooms/owned/:classroomId", handler.GetOwnedClassroom)
		route := fmt.Sprintf("/api/classrooms/owned/%d", 1)

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		remoteClassroom, err := query.Classroom.WithContext(context.Background()).Where(query.Classroom.OwnerID.Eq(1)).First()
		assert.NoError(t, err)
		assert.Equal(t, classroom.Name, remoteClassroom.Name)
		assert.Equal(t, classroom.Description, remoteClassroom.Description)
		assert.Equal(t, classroom.GroupID, remoteClassroom.GroupID)
		assert.Equal(t, classroom.GroupAccessTokenID, remoteClassroom.GroupAccessTokenID)
		assert.Equal(t, classroom.GroupAccessToken, remoteClassroom.GroupAccessToken)
	})
}
