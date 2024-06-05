package default_controller

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"

	"github.com/gofiber/fiber/v2"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestOwnedClassroomAssignmentProjectMiddleware(t *testing.T) {
	testDB := db_tests.NewTestDB(t)

	// Seeding data

	user := factory.User()
	testDB.InsertUser(&user)

	classroom := factory.Classroom()
	testDB.InsertClassroom(&classroom)

	userClassroom := factory.UserClassroom(user.ID, classroom.ID)
	testDB.InsertUserClassroom(&userClassroom)

	assignment := factory.Assignment(classroom.ID)
	testDB.InsertAssignment(&assignment)

	team := factory.Team(classroom.ID)
	testDB.InsertTeam(&team)

	project := factory.AssignmentProject(assignment.ID, team.ID)
	testDB.InsertAssignmentProjects(&project)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetUserID(1)
		return c.Next()
	})

	mailRepo := mailRepoMock.NewMockRepository(t)
	handler := NewApiController(mailRepo)
	app.Use("/api", handler.OwnedClassroomAssignmentProjectMiddleware)

	t.Run("ValidAssignmentProjectMiddlewareCall", func(t *testing.T) {
		route := fmt.Sprintf("/api/classrooms/%s/assignments/%s/projects/%s", classroom.ID, assignment.ID, project.ID)
		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("InvalidParametersMiddlewareCall", func(t *testing.T) {
		route := fmt.Sprintf("/api/classrooms/%s/assignments/%s/projects/%s", classroom.ID, uuid.New(), uuid.New())
		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)
		assert.Error(t, err)
		assert.NotEqual(t, fiber.StatusOK, resp.StatusCode)
	})
}
