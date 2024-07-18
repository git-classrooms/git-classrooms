package api

import (
	"context"
	"fmt"
	"gitlab.hs-flensburg.de/gitlab-classroom/config"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestDeleteClassroomInvitation(t *testing.T) {
	testDb := db_tests.NewTestDB(t)

	user1 := database.User{
		ID:             1,
		GitlabUsername: "user1",
		GitlabEmail:    "user1",
	}
	testDb.InsertUser(&user1)

	classroom := database.Classroom{
		ID:      uuid.New(),
		OwnerID: user1.ID,
	}
	testDb.InsertClassroom(&classroom)

	invitation := database.ClassroomInvitation{
		ID:          uuid.New(),
		Status:      database.ClassroomInvitationPending,
		ClassroomID: classroom.ID,
	}
	testDb.InsertInvitation(&invitation)

	app := fiber.New()
	mailRepo := mailRepoMock.NewMockRepository(t)
	handler := NewApiV2Controller(mailRepo, config.ApplicationConfig{})
	app.Delete("/classrooms/:classroomId/invitations/:invitationId", handler.RevokeClassroomInvitation)

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
		testDb.SaveInvitation(&invitation)

		req := httptest.NewRequest("DELETE", targetRoute, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("could not perform request: %s", err.Error())
		}

		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})

	t.Run("Revoke Classroom Invitation - Already Accepted", func(t *testing.T) {
		invitation.Status = database.ClassroomInvitationAccepted
		testDb.SaveInvitation(&invitation)

		req := httptest.NewRequest("DELETE", targetRoute, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("could not perform request: %s", err.Error())
		}

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Revoke Classroom Invitation - Already Rejected", func(t *testing.T) {
		invitation.Status = database.ClassroomInvitationRejected
		testDb.SaveInvitation(&invitation)

		req := httptest.NewRequest("DELETE", targetRoute, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("could not perform request: %s", err.Error())
		}

		assert.Equal(t, fiber.StatusAccepted, resp.StatusCode)
	})

	t.Run("Revoke Classroom Invitation", func(t *testing.T) {
		invitation.Status = database.ClassroomInvitationPending
		testDb.SaveInvitation(&invitation)

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
