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

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestOwnedClassroomMiddleware(t *testing.T) {
	testDB := db_tests.NewTestDB(t)
	// ------------ END OF DB SETUP -----------------

	user := factory.User()
	testDB.InsertUser(user)

	classroom := factory.Classroom()
	testDB.InsertClassroom(classroom)

	// ------------ END OF SEEDING DATA -----------------

	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetUserID(1)
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("OwnedClassroomMiddleware", func(t *testing.T) {
		app.Get("/api/classrooms/owned/:classroomId", handler.OwnedClassroomMiddleware, func(c *fiber.Ctx) error {
			classroom := fiberContext.Get(c).GetOwnedClassroom()
			return c.JSON(classroom)
		})

		route := fmt.Sprintf("/api/classrooms/owned/%s", classroom.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		retrievedClassroom, err := query.Classroom.WithContext(context.Background()).Where(query.Classroom.OwnerID.Eq(1)).First()
		assert.NoError(t, err)
		assert.Equal(t, classroom.Name, retrievedClassroom.Name)
	})
}
