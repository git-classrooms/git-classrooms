package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
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
	factory.UserClassroom(owner.ID, classroom.ID, database.Owner)

	// Setup fiber
	app, _, _ := setupApp(t, owner)

	t.Run("return a classroom by id", func(t *testing.T) {
		// prepare request
		route := fmt.Sprintf("/api/v1/classrooms/%s", classroom.ID)
		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Handle response
		var classroomResponse *UserClassroomResponse

		err = json.NewDecoder(resp.Body).Decode(&classroomResponse)
		assert.NoError(t, err)

		assert.Equal(t, classroomResponse.UserID, owner.ID)
		assert.Equal(t, classroomResponse.Role, database.Owner)
	})
}
