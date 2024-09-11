package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetClassroomTeamMembers(t *testing.T) {
	restoreDatabase(t)

	db, err := gorm.Open(postgres.Open(integrationTest.dbURL))
	if err != nil {
		t.Fatal(err)
	}

	query.SetDefault(db)

	owner := factory.User()
	classroom := factory.Classroom(owner.ID)
	factory.UserClassroom(owner.ID, classroom.ID, database.Owner)

	member := factory.User()

	members := []*database.UserClassrooms{
		factory.UserClassroom(member.ID, classroom.ID, database.Student),
	}

	team := factory.Team(classroom.ID, members)

	// ------------ END OF SEEDING DATA -----------------

	app, _, _ := setupApp(t, member)

	t.Run("TestGetClassroomTeamMembers", func(t *testing.T) {
		route := fmt.Sprintf("/api/v2/classrooms/%s/teams/%s/members", classroom.ID.String(), team.ID)

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var userClassroomResponse []*UserClassroomResponse

		err = json.NewDecoder(resp.Body).Decode(&userClassroomResponse)
		assert.NoError(t, err)

		assert.Len(t, userClassroomResponse, 1)

		assert.Equal(t, userClassroomResponse[0].TeamID, team.ID)
		assert.Equal(t, userClassroomResponse[0].UserID, member.ID)
	})
}
