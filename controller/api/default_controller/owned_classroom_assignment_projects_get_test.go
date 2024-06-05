package default_controller

import (
	"fmt"
	"net/http/httptest"
	"testing"

	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestGetOwnedClassroomAssignmentProjects(t *testing.T) {
	testDB := db_tests.NewTestDB(t)

	// Create test user
	user := factory.User()
	testDB.InsertUser(&user)

	classroom := factory.Classroom()
	testDB.InsertClassroom(&classroom)

	assignment := factory.Assignment(classroom.ID)
	testDB.InsertAssignment(&assignment)

	team := factory.Team(classroom.ID)
	testDB.InsertTeam(&team)

	factory.AssignmentProject(assignment.ID, team.ID)

	// Session and context setup
	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetOwnedClassroom(&classroom)
		ctx.SetOwnedClassroomAssignment(&assignment)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("GetOwnedClassroomAssignmentProjects", func(t *testing.T) {
		app.Get("/api/classrooms/owned/:classroomId/assignments/:assignmentId/projects", handler.GetOwnedClassroomAssignmentProjects)
		route := fmt.Sprintf("/api/classrooms/owned/%s/assignments/%s/projects", classroom.ID, assignment.ID)

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

	})
}
