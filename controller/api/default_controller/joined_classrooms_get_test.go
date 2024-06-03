package default_controller

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

func TestGetJoinedClassrooms(t *testing.T) {
	testDB := db_tests.NewTestDB(t)

	owner := factory.User()
	testDB.InsertUser(owner)

	classroom := factory.Classroom()
	testDB.InsertClassroom(classroom)

	joinedClassrooms := []*database.Classroom{
		{
			Name:               "Joined Classroom One",
			OwnerID:            1,
			Description:        "Description One",
			GroupID:            15,
			GroupAccessTokenID: 35,
			GroupAccessToken:   "token35",
		},
		{
			Name:               "Joined Classroom Two",
			OwnerID:            1,
			Description:        "Description Two",
			GroupID:            25,
			GroupAccessTokenID: 45,
			GroupAccessToken:   "token45",
		},
	}

	for _, classroom := range joinedClassrooms {
		testDB.InsertClassroom(classroom)
	}

	// ------------ END OF SEEDING DATA -----------------

	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetUserID(1)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("GetJoinedClassrooms", func(t *testing.T) {
		app.Get("/api/classrooms/joined", handler.GetJoinedClassrooms)
		route := fmt.Sprintf("/api/classrooms/joined")

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		type ClassRoomResponse struct {
			ID          uuid.UUID `json:"id"`
			Name        string    `json:"name"`
			OwnerID     int       `json:"ownerId"`
			Description string    `json:"description"`
			GroupID     int       `json:"groupId"`
		}

		var classrooms []*ClassRoomResponse

		err = json.NewDecoder(resp.Body).Decode(&classrooms)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		assert.Len(t, classrooms, len(joinedClassrooms))
		for i, classroom := range classrooms {
			assert.Equal(t, joinedClassrooms[i].ID, classroom.ID)
			assert.Equal(t, joinedClassrooms[i].Name, classroom.Name)
			assert.Equal(t, joinedClassrooms[i].OwnerID, classroom.OwnerID)
			assert.Equal(t, joinedClassrooms[i].Description, classroom.Description)
			assert.Equal(t, joinedClassrooms[i].GroupID, classroom.GroupID)
		}
	})
}
