package default_controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestGetJoinedClassrooms(t *testing.T) {
	// --------------- DB SETUP -----------------
	t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "false")
	pg, err := tests.StartPostgres()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err = pg.Restore(context.Background())
		if err != nil {
			t.Fatal(err)
		}
	})
	dbURL, err := pg.ConnectionString(context.Background())
	db, err := gorm.Open(postgresDriver.Open(dbURL), &gorm.Config{})
	if err != nil {
		t.Fatalf("could not connect to database: %s", err.Error())
	}
	err = utils.MigrateDatabase(db)
	if err != nil {
		t.Fatalf("could not migrate database: %s", err.Error())
	}
	err = pg.Snapshot(context.Background(), postgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		t.Fatal(err)
	}
	query.SetDefault(db)

	// ------------ END OF DB SETUP -----------------

	owner := &database.User{ID: 1}
	err = query.User.WithContext(context.Background()).Create(owner)
	if err != nil {
		t.Fatalf("could not create test owner: %s", err.Error())
	}

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
		err = query.Classroom.WithContext(context.Background()).Create(classroom)
		if err != nil {
			t.Fatalf("could not create test classroom: %s", err.Error())
		}
	}

	// ------------ END OF SEEDING DATA -----------------

	session.InitSessionStore(dbURL)
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
		route := fmt.Sprintf("/api/classrooms/joined", nil)

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