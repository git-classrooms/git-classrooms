package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
)

func TestGetClassroomTeam(t *testing.T) {
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

	dueDate := time.Now().Add(1 * time.Hour)

	members := []*database.UserClassrooms{
		factory.UserClassroom(member.ID, classroom.ID, database.Student),
	}

	factory.Assignment(classroom.ID, &dueDate, false)
	team := factory.Team(classroom.ID, members)

	// ------------ END OF SEEDING DATA -----------------

	app, _, _ := setupApp(t, member)

	t.Run("TestGetClassroomTeam", func(t *testing.T) {
		route := fmt.Sprintf("/api/v1/classrooms/%s/teams/%s", classroom.ID.String(), team.ID)

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var teamResponse *TeamResponse

		err = json.NewDecoder(resp.Body).Decode(&teamResponse)
		assert.NoError(t, err)
		assert.Equal(t, team.ID, teamResponse.ID)
	})
}
