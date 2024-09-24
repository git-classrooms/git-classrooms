package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
)

func TestGetClassroomMembers(t *testing.T) {
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

	t.Run("GetClassroomMembers", func(t *testing.T) {
		route := fmt.Sprintf("/api/v1/classrooms/%s/members", classroom.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var membersResponse []*UserClassroomResponse

		err = json.NewDecoder(resp.Body).Decode(&membersResponse)
		assert.NoError(t, err)

		assert.Len(t, membersResponse, 2)

		assert.Equal(t, owner.ID, membersResponse[0].User.ID)
		assert.Equal(t, member.ID, membersResponse[1].User.ID)
	})
}
