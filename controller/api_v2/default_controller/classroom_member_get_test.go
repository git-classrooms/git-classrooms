package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetClassroomMember(t *testing.T) {
	restoreDatabase(t)

	db, err := gorm.Open(postgres.Open(integrationTest.dbURL))
	if err != nil {
		t.Fatal(err)
	}

	query.SetDefault(db)

	owner := factory.User()
	classroom := factory.Classroom(owner.ID)
	factory.UserClassroom(owner.ID, classroom.ID, database.Owner)

	member := factory.User()
	factory.UserClassroom(member.ID, classroom.ID, database.Student)

	// ------------ END OF SEEDING DATA -----------------

	app, _, _ := setupApp(t, owner)

	t.Run("GetClassroomMember", func(t *testing.T) {
		route := fmt.Sprintf("/api/v2/classrooms/%s/members/%d", classroom.ID.String(), member.ID)

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)

		var classroomResponse *UserClassroomResponse

		err = json.NewDecoder(resp.Body).Decode(&classroomResponse)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		assert.Equal(t, member.ID, classroomResponse.User.ID)
	})
}
