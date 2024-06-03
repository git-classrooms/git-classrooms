package default_controller

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestGetOwnedClassroomAssignmentProject(t *testing.T) {

	testDB := db_tests.NewTestDB(t)


	// Create test user
	owner := factory.User()
	testDB.InsertUser(owner )

	classroom := factory.Classroom()
	testDB.InsertClassroom(classroom)

	assignment := factory.Assignment(classroom.ID)
	testDB.InsertAssignment(assignment)

	team := factory.Team(classroom.ID)
	testDB.InsertTeam(team)

	project := factory.AssignmentProject(assignment.ID, team.ID)
	testDB.InsertAssignmentProjects(project)


	// ------------ END OF SEEDING DATA -----------------

	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetOwnedClassroom(classroom)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(owner.ID)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("GetOwnedClassroomAssignmentProject", func(t *testing.T) {
		app.Get("/classrooms/owned/:classroomId/assignments/:assignmentId/projects/:projectId", handler.GetOwnedClassroomAssignment)
		route := fmt.Sprintf("/api/classrooms/owned/%s/assignments/%s/projects/%s", classroom.ID.String(), assignment.ID.String(), project.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)

		type ClassroomAssignmentProjectResponse struct {
			ID                 uuid.UUID `json:"id"`
			CreatedAt          time.Time `json:"createdAt"`
			UpdatedAt          time.Time `json:"updatedAt"`
			AssignmentID       uuid.UUID `json:"assignmentId"`
			UserID             int       `json:"userId"`
			AssignmentAccepted bool      `json:"assignmentAccepted"`
			ProjectID          int       `json:"projectId"`
		}

		var classroomAssignmentProject *ClassroomAssignmentProjectResponse

		err = json.NewDecoder(resp.Body).Decode(&classroomAssignmentProject)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		assert.Equal(t, project.ID, classroomAssignmentProject.ID)
		assert.Equal(t, project.AssignmentID, classroomAssignmentProject.AssignmentID)
		// FIXME: NEEDS fix
		assert.Equal(t, project.Team.ID, classroomAssignmentProject.UserID)
		assert.Equal(t, project.AssignmentAccepted, classroomAssignmentProject.AssignmentAccepted)
		assert.Equal(t, project.ProjectID, classroomAssignmentProject.ProjectID)
	})
}
