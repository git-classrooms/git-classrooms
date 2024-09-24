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

func TestGetProjects(t *testing.T) {
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

	assignment := factory.Assignment(classroom.ID, &dueDate, false)
	team := factory.Team(classroom.ID, members)
	assignmentProject := factory.AssignmentProject(assignment.ID, team.ID)

	// ------------ END OF SEEDING DATA -----------------

	app, _, _ := setupApp(t, member)

	t.Run("TestGetProjects", func(t *testing.T) {
		route := fmt.Sprintf("/api/v1/classrooms/%s/projects", classroom.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var projectsResponse []*ProjectResponse

		err = json.NewDecoder(resp.Body).Decode(&projectsResponse)
		assert.NoError(t, err)

		assert.Len(t, projectsResponse, 1)

		assert.Equal(t, assignmentProject.ID, projectsResponse[0].ID)
	})
}
