package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/config"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetClassroom(t *testing.T) {
	restoreDatabase(t)

	db, err := gorm.Open(postgres.Open(integrationTest.dbURL))
	if err != nil {
		t.Fatal(err)
	}

	query.SetDefault(db)

	owner := factory.User()
	classroom := factory.Classroom(owner.ID)
	userClassroom := factory.UserClassroom(owner.ID, classroom.ID, database.Owner)

	// Setup fiber
	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := context.Get(c)
		ctx.SetUserID(owner.ID)
		ctx.SetUserClassroom(userClassroom)
		return c.Next()

	})

	mailRepo := mailRepoMock.NewMockRepository(t)
	handler := NewApiV2Controller(mailRepo, config.ApplicationConfig{})
	app.Get("/api/v2/classrooms/:classroomId", handler.GetClassroom)

	t.Run("return a classroom by id", func(t *testing.T) {
		// prepare request
		route := fmt.Sprintf("/api/v2/classrooms/%s", classroom.ID)
		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)

		// Handle response
		var classroomResponse *UserClassroomResponse

		err = json.NewDecoder(resp.Body).Decode(&classroomResponse)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		assert.Equal(t, classroomResponse.UserID, owner.ID)
		assert.Equal(t, classroomResponse.Role, database.Owner)
	})
}
