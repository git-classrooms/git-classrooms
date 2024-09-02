package api

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
)

func TestDeleteClassroomInvitation(t *testing.T) {
	restoreDatabase(t)

	owner := factory.User() // id 0
	classroom := factory.Classroom(owner.ID)
	factory.UserClassroom(owner.ID, classroom.ID, database.Owner)
	invitation := factory.Invitation(classroom.ID)

	app := setupApp(t, owner, nil)

	targetRoute := fmt.Sprintf("/api/v2/classrooms/%s/invitations/%s", invitation.ClassroomID.String(), invitation.ID.String())

	t.Run("Revoke Classroom Invitation - Not Found", func(t *testing.T) {
		newTarget := fmt.Sprintf("/api/v2/classrooms/%s/invitations/%s", uuid.New(), uuid.New())
		req := httptest.NewRequest("DELETE", newTarget, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("could not perform request: %s", err.Error())
		}

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Revoke Classroom Invitation - Already Revoked", func(t *testing.T) {
		invitation.Status = database.ClassroomInvitationRevoked
		SaveInvitation(t, invitation)

		req := httptest.NewRequest("DELETE", targetRoute, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("could not perform request: %s", err.Error())
		}

		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})

	t.Run("Revoke Classroom Invitation - Already Accepted", func(t *testing.T) {
		invitation.Status = database.ClassroomInvitationAccepted
		SaveInvitation(t, invitation)

		req := httptest.NewRequest("DELETE", targetRoute, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("could not perform request: %s", err.Error())
		}

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Revoke Classroom Invitation - Already Rejected", func(t *testing.T) {
		invitation.Status = database.ClassroomInvitationRejected
		SaveInvitation(t, invitation)

		req := httptest.NewRequest("DELETE", targetRoute, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("could not perform request: %s", err.Error())
		}

		assert.Equal(t, fiber.StatusAccepted, resp.StatusCode)
	})

	t.Run("Revoke Classroom Invitation", func(t *testing.T) {
		invitation.Status = database.ClassroomInvitationPending
		SaveInvitation(t, invitation)

		req := httptest.NewRequest("DELETE", targetRoute, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("could not perform request: %s", err.Error())
		}

		assert.Equal(t, fiber.StatusAccepted, resp.StatusCode)

		dbInvitation, err := query.ClassroomInvitation.WithContext(context.Background()).Where(query.ClassroomInvitation.ID.Eq(invitation.ID)).First()
		assert.NoError(t, err)
		assert.Equal(t, database.ClassroomInvitationRevoked, dbInvitation.Status)
	})

}

func SaveInvitation(t *testing.T, invitation *database.ClassroomInvitation) {
	err := query.ClassroomInvitation.WithContext(context.Background()).Save(invitation)
	if err != nil {
		t.Fatalf("could not update invitation: %s", err.Error())
	}
}
