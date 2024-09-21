package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetClassroomInvitation(t *testing.T) {
	restoreDatabase(t)

	db, err := gorm.Open(postgres.Open(integrationTest.dbURL))
	if err != nil {
		t.Fatal(err)
	}

	query.SetDefault(db)

	user := factory.User()
	classroom := factory.Classroom(user.ID)
	factory.UserClassroom(user.ID, classroom.ID, database.Owner)
	invitation := factory.Invitation(classroom.ID)

	app, _, _ := setupApp(t, user)

	route := fmt.Sprintf("/api/v1/classrooms/%s/invitations/%s", classroom.ID, invitation.ID)

	req := httptest.NewRequest("GET", route, nil)
	resp, err := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.NoError(t, err)

	t.Log(resp.Body)

	var invitationResponse *database.ClassroomInvitation

	err = json.NewDecoder(resp.Body).Decode(&invitation)
	t.Log(invitationResponse)
	assert.NoError(t, err)

	// t.Log(invitation.Status)
}
