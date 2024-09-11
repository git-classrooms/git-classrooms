package api

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestPutClassroom(t *testing.T) {
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
	app, gitlabRepo, _ := setupApp(t, owner)

	t.Run("TestPutClassroom", func(t *testing.T) {
		// prepare request
		route := fmt.Sprintf("/api/v2/classrooms/%s", classroom.ID)
		requestBody := updateClassroomRequest{
			Name:        "new name",
			Description: "new Description",
		}

		gitlabRepo.
			EXPECT().
			ChangeGroupName(classroom.GroupID, requestBody.Name).
			Return(&model.Group{Name: requestBody.Name}, nil)

		gitlabRepo.
			EXPECT().
			ChangeGroupDescription(classroom.GroupID, requestBody.Description).
			Return(&model.Group{Description: requestBody.Description}, nil)

		req := newPutJsonRequest(route, requestBody)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusAccepted, resp.StatusCode)

		// Handle response
		var classroomResponse *database.Classroom

		err = json.NewDecoder(resp.Body).Decode(&classroomResponse)
		assert.NoError(t, err)

		assert.Equal(t, classroomResponse.Name, requestBody.Name)
		assert.Equal(t, classroomResponse.Description, requestBody.Description)
	})
}
