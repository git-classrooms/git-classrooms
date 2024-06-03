package default_controller

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestInviteToClassroom(t *testing.T) {
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
		ctx.SetOwnedClassroom(classroom)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("InviteToClassroom", func(t *testing.T) {
		app.Post("/api/classrooms/:classroomId/invitations", handler.InviteToClassroom)
		route := fmt.Sprintf("/api/classrooms/%d/invitations", classroom.ID)

		req := httptest.NewRequest("POST", route, strings.NewReader(`{"memberEmails":["test1@example.com", "test2@example.com"]}`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		invitations, err := query.ClassroomInvitation.
			WithContext(context.Background()).
			Where(query.ClassroomInvitation.ClassroomID.Eq(classroom.ID)).
			Find()
		assert.NoError(t, err)
		assert.Len(t, invitations, 2)
		assert.Equal(t, "test1@example.com", invitations[0].Email)
		assert.Equal(t, "test2@example.com", invitations[1].Email)
	})
}
