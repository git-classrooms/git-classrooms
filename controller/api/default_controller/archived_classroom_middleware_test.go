package api

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"gitlab.hs-flensburg.de/gitlab-classroom/config"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func TestArchivedClassroomMiddleware(t *testing.T) {
	userClassroom := database.UserClassrooms{
		Classroom: database.Classroom{
			ID:       uuid.New(),
			Archived: true,
		},
	}

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := context.Get(c)
		ctx.SetUserClassroom(&userClassroom)
		return c.Next()
	})

	mailRepo := mailRepoMock.NewMockRepository(t)
	handler := NewAPIV1Controller(mailRepo, config.ApplicationConfig{})
	app.Use("/api/v1/classrooms/:classroomId", handler.ArchivedMiddleware)

	targetRoute := fmt.Sprintf("/api/v1/classrooms/%s", userClassroom.Classroom.ID.String())

	t.Run("Forbiddes post for archived classroom", func(t *testing.T) {
		testForbiddenMethod(t, app, targetRoute, fiber.MethodPost)
	})

	t.Run("Forbiddes put for archived classroom", func(t *testing.T) {
		testForbiddenMethod(t, app, targetRoute, fiber.MethodPut)
	})

	t.Run("Forbiddes patch for archived classroom", func(t *testing.T) {
		testForbiddenMethod(t, app, targetRoute, fiber.MethodPatch)
	})

	t.Run("Forbiddes delete for archived classroom", func(t *testing.T) {
		testForbiddenMethod(t, app, targetRoute, fiber.MethodDelete)
	})

	t.Run("Allows get for archived classroom", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, targetRoute, nil)
		resp, err := app.Test(req)
		assert.NotEqual(t, fiber.StatusForbidden, resp.StatusCode)
		assert.NoError(t, err)
		defer resp.Body.Close()
	})
}

func testForbiddenMethod(t *testing.T, app *fiber.App, targetRoute string, method string) {
	req := httptest.NewRequest(method, targetRoute, nil)
	resp, err := app.Test(req)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	assert.NoError(t, err)
	defer resp.Body.Close()
}
