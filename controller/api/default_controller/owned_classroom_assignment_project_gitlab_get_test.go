package default_controller

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestGetOwnedClassroomAssignmentProjectGitlab(t *testing.T) {
	testDB := db_tests.NewTestDB(t)

	// Create test user
	user := factory.User()
	testDB.InsertUser(user)

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
		ctx.SetOwnedClassroomAssignmentProject(project)
		ctx.SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("GetOwnedClassroomAssignmentProjectGitlab", func(t *testing.T) {
		app.Get("/api/classrooms/owned/:classroomId/assignments/:assignmentId/gitlab", handler.GetOwnedClassroomAssignmentProject)
		route := fmt.Sprintf("/api/classrooms/owned/%s/assignments/%s/gitlab", classroom.ID, assignment.ID)

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusFound, resp.StatusCode)
		assert.NoError(t, err)
	})
}
