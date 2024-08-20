package api

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	test_db "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestDeleteClassroomInvitation(t *testing.T) {
	restoreDatabase(t)

	db, err := gorm.Open(postgres.Open(integrationTest.dbURL))
	if err != nil {
		t.Fatal(err)
	}

	query.SetDefault(db)
	user := factory.User()
	classroom := factory.Classroom(user.ID)
	invitation := factory.Invitation(classroom.ID)

	app := setupApp(t, user, nil)

	targetRoute := fmt.Sprintf("/classrooms/%s/invitations/%s", invitation.ClassroomID, invitation.ID)

	t.Run("Revoke Classroom Invitation - Not Found", func(t *testing.T) {
		newTarget := fmt.Sprintf("/classrooms/%s/invitations/%s", uuid.New(), uuid.New())
		req := httptest.NewRequest("DELETE", newTarget, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("could not perform request: %s", err.Error())
		}

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Revoke Classroom Invitation - Already Revoked", func(t *testing.T) {
		invitation.Status = database.ClassroomInvitationRevoked
		test_db.SaveInvitation(t, invitation)

		req := httptest.NewRequest("DELETE", targetRoute, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("could not perform request: %s", err.Error())
		}

		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})

	t.Run("Revoke Classroom Invitation - Already Accepted", func(t *testing.T) {
		invitation.Status = database.ClassroomInvitationAccepted
		test_db.SaveInvitation(t, invitation)

		req := httptest.NewRequest("DELETE", targetRoute, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("could not perform request: %s", err.Error())
		}

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Revoke Classroom Invitation - Already Rejected", func(t *testing.T) {
		invitation.Status = database.ClassroomInvitationRejected
		test_db.SaveInvitation(t, invitation)

		req := httptest.NewRequest("DELETE", targetRoute, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("could not perform request: %s", err.Error())
		}

		assert.Equal(t, fiber.StatusAccepted, resp.StatusCode)
	})

	t.Run("Revoke Classroom Invitation", func(t *testing.T) {
		invitation.Status = database.ClassroomInvitationPending
		test_db.SaveInvitation(t, invitation)

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
