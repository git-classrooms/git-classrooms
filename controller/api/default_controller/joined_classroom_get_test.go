package default_controller

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestGetJoinedClassroom(t *testing.T) {
	testDB := db_tests.NewTestDB(t)

	owner := database.User{ID: 1, GitlabEmail: "owner@example.com"}
	testDB.InsertUser(&owner)

	member := database.User{ID: 2, GitlabEmail: "member@example.com"}
	testDB.InsertUser(&member)

	classroom := factory.Classroom(map[string]any{"OwnerID": owner.ID})
	testDB.InsertClassroom(&classroom)

	userClassroom := &database.UserClassrooms{
		UserID:      member.ID,
		ClassroomID: classroom.ID,
		Role:        database.Student,
	}
	testDB.InsertUserClassroom(userClassroom)

	// ------------ END OF SEEDING DATA -----------------

	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetOwnedClassroom(&classroom)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(member.ID)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("GetJoinedClassroom", func(t *testing.T) {
		app.Get("/api/classrooms/joined/:classroomId", handler.GetJoinedClassroom)
		route := fmt.Sprintf("/api/classrooms/joined/%s", classroom.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)

		type ClassRoomResponse struct {
			ID          uuid.UUID `json:"id"`
			Name        string    `json:"name"`
			OwnerID     int       `json:"ownerId"`
			Description string    `json:"description"`
			GroupID     int       `json:"groupId"`
		}

		var classRoom *ClassRoomResponse

		err = json.NewDecoder(resp.Body).Decode(&classRoom)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		assert.Equal(t, classroom.ID, classRoom.ID)
		assert.Equal(t, classroom.Name, classRoom.Name)
		assert.Equal(t, classroom.OwnerID, classRoom.OwnerID)
		assert.Equal(t, classroom.Description, classRoom.Description)
		assert.Equal(t, classroom.GroupID, classRoom.GroupID)
	})
}
