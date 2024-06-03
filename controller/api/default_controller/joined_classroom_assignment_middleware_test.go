package default_controller

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJoinedClassroomAssignmentMiddleware(t *testing.T) {
	testDB := db_tests.NewTestDB(t)

	// Create test user
	user := database.User{ID: 1, GitlabEmail: "test@example.com", Name: "Test User"}
	testDB.InsertUser(&user)

	// Create test classroom
	classroom := database.Classroom{
		ID:                 uuid.UUID{},
		Name:               "Test classroom",
		OwnerID:            1,
		Description:        "Classroom description",
		GroupID:            1,
		GroupAccessTokenID: 20,
		GroupAccessToken:   "token",
	}
	testDB.InsertClassroom(&classroom)

	assignment := database.Assignment{
		ID:                uuid.UUID{},
		ClassroomID:       classroom.ID,
		TemplateProjectID: 1234,
		Name:              "Test Assignment",
		Description:       "Test Assignment Description",
	}
	testDB.InsertAssignment(&assignment)

	team := database.Team{
		ID:          uuid.UUID{},
		ClassroomID: classroom.ID,
	}
	testDB.InsertTeam(&team)

	// ------------ END OF SEEDING DATA -----------------

	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("JoinedClassroomAssignmentMiddleware", func(t *testing.T) {
		app.Use("/api/classrooms/:classroomId/assignments/:assignmentId", handler.JoinedClassroomAssignmentMiddleware)
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
