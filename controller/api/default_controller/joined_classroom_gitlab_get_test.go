package default_controller

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestGetJoinedClassroomGitlab(t *testing.T) {
	testDB := db_tests.NewTestDB(t)

	user := factory.User()
	testDB.InsertUser(user)

	classroom := factory.Classroom()
	testDB.InsertClassroom(classroom)

	// ------------ END OF SEEDING DATA -----------------

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

	t.Run("GetJoinedClassroomGitlab", func(t *testing.T) {
		app.Get("/api/v1//classrooms/joined/:classroomId/gitlab", handler.RedirectGroupGitlab)
		route := fmt.Sprintf("/api/classrooms/joined/%d/gitlab", 1)

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusFound, resp.StatusCode) // Check for redirection status
		assert.NoError(t, err)

		location := resp.Header.Get("Location")
		assert.NotEmpty(t, location)
		assert.Contains(t, location, "group") // Adjust this to match the actual expected URL structure
	})
}
