package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
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

	assignment := factory.Assignment(classroom.ID)
	team := factory.Team(classroom.ID)
	project := factory.AssignmentProject(assignment.ID, team.ID)

	// setup app
	mailRepo := mailRepoMock.NewMockRepository(t)
	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetOwnedClassroom(&classroom)
		ctx.SetUserClassroom(&userClassroom)
		ctx.SetAssignment(&assignment)

		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(owner.ID)
		s.Save()
		return c.Next()
	})

	handler := NewApiV2Controller(mailRepo)

	t.Run("GetOwnedClassroomAssignmentProjects", func(t *testing.T) {
		app.Get("/api/v2/classrooms/:classroomId/assignments/:assignmentId/projects", handler.GetClassroomAssignmentProjects)
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
		assert.Equal(t, project.AssignmentAccepted, projectResponse.AssignmentAccepted)
		assert.Equal(t, project.ProjectID, projectResponse.ProjectID)
	})
}
