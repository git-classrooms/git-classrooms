package default_controller

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestGetOwnedClassroomInvitations(t *testing.T) {
	testDB := db_tests.NewTestDB(t)

	user := factory.User()
	testDB.InsertUser(&user)

	classroom := factory.Classroom()
	testDB.InsertClassroom(&classroom)

	invitation := factory.Invitation(classroom.ID)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetOwnedClassroom(&classroom)

		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(nil)

	t.Run("GetOwnedClassroomInvitations", func(t *testing.T) {
		app.Get("/api/classrooms/owned/:classroomId/invitations", handler.GetOwnedClassroomInvitations)
		route := fmt.Sprintf("/api/classrooms/owned/%s/invitations", classroom.ID)

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		invitations, err := query.ClassroomInvitation.
			WithContext(context.Background()).
			Where(query.ClassroomInvitation.ClassroomID.Eq(classroom.ID)).
			Find()

		assert.NoError(t, err)
		assert.NotEmpty(t, invitations)
		assert.Equal(t, invitation.Email, invitations[0].Email)
	})
}
