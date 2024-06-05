package default_controller

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJoinedClassroomAssignmentMiddleware(t *testing.T) {
	testDB := db_tests.NewTestDB(t)

	// Create test user
	user := factory.User()
	testDB.InsertUser(&user)

	classroom := factory.Classroom()
	testDB.InsertClassroom(&classroom)

	assignment := factory.Assignment(classroom.ID)
	testDB.InsertAssignment(&assignment)

	team := database.Team{
		ID:          uuid.UUID{},
		ClassroomID: classroom.ID,
	}
	testDB.InsertTeam(&team)

	project := factory.AssignmentProject(assignment.ID, team.ID)
	testDB.InsertAssignmentProjects(&project)

	assignment.Projects = append(assignment.Projects, &project)

	// ------------ END OF SEEDING DATA -----------------

	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(user.ID)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("JoinedClassroomAssignmentMiddleware", func(t *testing.T) {
		app.Use("/api/v1/classrooms/:classroomId/assignments/:assignmentId", handler.JoinedClassroomAssignmentMiddleware)
		route := fmt.Sprintf("/api/v1/classrooms/%s/assignments/%s", classroom.ID.String(), assignment.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		assignmentProject, err := query.AssignmentProjects.WithContext(context.Background()).
			Where(query.AssignmentProjects.AssignmentID.Eq(assignment.ID)).
			First()
		assert.NoError(t, err)
		assert.Equal(t, assignment.ID, assignmentProject.AssignmentID)
		assert.Equal(t, team.ID, assignmentProject.Team.ID)
	})
}
