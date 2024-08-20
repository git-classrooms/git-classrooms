package api

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	//"gitlab.hs-flensburg.de/gitlab-classroom/config"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	// mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	//"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetClassrooms(t *testing.T) {
	restoreDatabase(t)

	//mailRepo := mailRepoMock.NewMockRepository(t)

	db, err := gorm.Open(postgres.Open(integrationTest.dbURL))
	if err != nil {
		t.Fatal(err)
	}

	query.SetDefault(db)

	owner := factory.User()
	classroom := factory.Classroom(owner.ID)
	moderator := factory.User()
	student := factory.User()

	factory.UserClassroom(owner.ID, classroom.ID, database.Owner)
	factory.UserClassroom(moderator.ID, classroom.ID, database.Moderator)
	factory.UserClassroom(student.ID, classroom.ID, database.Student)

	t.Run("return all classrooms where the user is the owner", func(t *testing.T) {
		app := setupApp(t, owner, nil)
		// prepare request
		route := "/api/v2/classrooms?filter=owned"
		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)

		// Handle response
		var classroomsResponse []*UserClassroomResponse

		err = json.NewDecoder(resp.Body).Decode(&classroomsResponse)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		// Check response
		assert.Len(t, classroomsResponse, 1)

		classroomResponse := classroomsResponse[0]

		assert.Equal(t, classroomResponse.UserID, owner.ID)
		assert.Equal(t, classroomResponse.Classroom.OwnerID, classroom.OwnerID)
		assert.Equal(t, classroomResponse.Classroom.ID, classroom.ID)
		assert.Equal(t, classroomResponse.Role, database.Owner)
	})

	t.Run("return all classrooms where the user is moderator", func(t *testing.T) {
		app := setupApp(t, moderator, nil)
		// prepare request
		route := "/api/v2/classrooms?filter=moderator"
		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)

		// Handle response
		var classroomsResponse []*UserClassroomResponse

		err = json.NewDecoder(resp.Body).Decode(&classroomsResponse)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		// Check response
		assert.Len(t, classroomsResponse, 1)

		classroomResponse := classroomsResponse[0]

		assert.Equal(t, classroomResponse.User.GitlabEmail, moderator.GitlabEmail)
		assert.Equal(t, classroomResponse.Classroom.OwnerID, classroom.OwnerID)
		assert.Equal(t, classroomResponse.Classroom.ID, classroom.ID)
		assert.Equal(t, classroomResponse.Role, database.Moderator)
	})

	t.Run("return all classrooms where the user is student", func(t *testing.T) {
		app := setupApp(t, student, nil)
		// prepare request
		route := "/api/v2/classrooms?filter=student"
		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)

		// Handle response
		var classroomsResponse []*UserClassroomResponse

		err = json.NewDecoder(resp.Body).Decode(&classroomsResponse)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		// Check response
		assert.Len(t, classroomsResponse, 1)

		classroomResponse := classroomsResponse[0]

		assert.Equal(t, classroomResponse.User.GitlabEmail, student.GitlabEmail)
		assert.Equal(t, classroomResponse.Classroom.OwnerID, classroom.OwnerID)
		assert.Equal(t, classroomResponse.Classroom.ID, classroom.ID)
		assert.Equal(t, classroomResponse.Role, database.Student)
	})
}
