package default_controller

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestOwnedClassroomAssignmentMiddleware(t *testing.T) {
	testDB := db_tests.NewTestDB(t)

	classroom := factory.Classroom()
	testDB.InsertClassroom(classroom)

	assignment := factory.Assignment(classroom.ID)
	testDB.InsertAssignment(assignment)

	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		fctx := fiberContext.Get(c)
		fctx.SetOwnedClassroom(classroom)

		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("OwnedClassroomAssignmentMiddleware", func(t *testing.T) {
		app.Get("/api/classrooms/:classroomId/assignments/:assignmentId", handler.OwnedClassroomAssignmentMiddleware)
		route := fmt.Sprintf("/api/classrooms/%s/assignments/%s", classroom.ID.String(), assignment.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		queriedAssignment, err := query.Assignment.WithContext(context.Background()).Where(query.Assignment.ID.Eq(assignment.ID)).First()
		assert.NoError(t, err)
		assert.Equal(t, assignment.Name, queriedAssignment.Name)
	})
}
