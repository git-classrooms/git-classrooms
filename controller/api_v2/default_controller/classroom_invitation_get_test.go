package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func TestGetOwnedClassroomInvitations(t *testing.T) {
	restoreDatabase(t)

	user := factory.User()
	classroom := factory.Classroom(user.ID)
	userClassroom := factory.UserClassroom(user.ID, classroom.ID, database.Owner)
	invitation := factory.Invitation(classroom.ID)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetUserClassroom(&userClassroom)

		return c.Next()
	})

	handler := NewApiV2Controller(nil)

	t.Run("GetClassroomInvitations", func(t *testing.T) {
		app.Get("/api/v2/classrooms/:classroomId/invitations", handler.GetClassroomInvitations)
		route := fmt.Sprintf("/api/v2/classrooms/%s/invitations", classroom.ID)

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		var invitations []*database.ClassroomInvitation

		err = json.NewDecoder(resp.Body).Decode(&invitations)
		assert.NoError(t, err)

		assert.Len(t, invitations, 1)
		assert.Equal(t, invitation.Email, invitations[0].Email)
	})
}
