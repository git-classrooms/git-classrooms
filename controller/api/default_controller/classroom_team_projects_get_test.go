package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
)

func TestGetClassroomTeamProjects(t *testing.T) {
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
	dueDate := time.Now().Add(1 * time.Hour)

	assignment := factory.Assignment(classroom.ID, &dueDate, false)
	project := factory.AssignmentProject(assignment.ID, team.ID)

	// ------------ END OF SEEDING DATA -----------------

	app, _, _ := setupApp(t, owner)

	t.Run("TestGetClassroomTeamMembers", func(t *testing.T) {
		route := fmt.Sprintf("/api/v1/classrooms/%s/teams/%s/projects", classroom.ID.String(), team.ID)

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var projectsResponse []*ProjectResponse

		err = json.NewDecoder(resp.Body).Decode(&projectsResponse)
		assert.NoError(t, err)

		assert.Len(t, projectsResponse, 1)

		assert.Equal(t, projectsResponse[0].Assignment.ID, assignment.ID)
		assert.Equal(t, projectsResponse[0].ID, project.ID)
	})
}
