package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetClassroomAssignmentProjects(t *testing.T) {
	restoreDatabase(t)

	db, err := gorm.Open(postgres.Open(integrationTest.dbURL))
	if err != nil {
		t.Fatal(err)
	}

	query.SetDefault(db)

	// Seed data
	owner := factory.User()
	classroom := factory.Classroom(owner.ID)
	userClassroom := factory.UserClassroom(owner.ID, classroom.ID, database.Owner)

	dueDate := time.Now().Add(1 * time.Hour)
	assignment := factory.Assignment(classroom.ID, &dueDate)
	team := factory.Team(classroom.ID, []*database.UserClassrooms{userClassroom})
	project := factory.AssignmentProject(assignment.ID, team.ID)

	// setup app
	app := setupApp(t, owner, nil)

	t.Run("GetOwnedClassroomAssignmentProjects", func(t *testing.T) {
		route := fmt.Sprintf("/api/v2/classrooms/%s/assignments/%s/projects", classroom.ID.String(), assignment.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)

		var projectsResponse []*ProjectResponse

		err = json.NewDecoder(resp.Body).Decode(&projectsResponse)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Len(t, projectsResponse, 1)

		projectResponse := projectsResponse[0]

		assert.Equal(t, project.ID, projectResponse.ID)
		assert.Equal(t, project.ProjectID, projectResponse.ProjectID)
	})
}
