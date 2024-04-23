package default_controller

import (
	"context"
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
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestGetOwnedClassroom(t *testing.T) {
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

	classroomsQuery := query.Classrooms
	testClassrooms := []*database.Classroom{
		{
			Name:                  "Test Classroom 1",
			Description:           "This is test classroom 1",
			OwnerID:               1,
			GroupID:               1,
			GroupAccessTokenID:    1,
			GroupAccessToken:      "test-token-1",
		},
		{
			Name:                  "Test Classroom 2",
			Description:           "This is test classroom 2",
			OwnerID:               1,
			GroupID:               2,
			GroupAccessTokenID:    2,
			GroupAccessToken:      "test-token-2",
		},
		{
			Name:                  "Test Classroom 3",
			Description:           "This is test classroom 3",
			OwnerID:               1,
			GroupID:               3,
			GroupAccessTokenID:    3,
			GroupAccessToken:      "test-token-3",
		},
	}

	for _, classroom := range testClassrooms {
		err = classroomQuery.WithContext(context.Background()).Create(testClassRoom)
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
		ctx.SetOwnedClassrooms(testClassrooms)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("GetOwnedClassrooms", func(t *testing.T) {
		app.Get("/api/classrooms/owned", handler.GetOwnedClassrooms)
		route := fmt.Sprintf("/api/classrooms/owned")

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		var classrooms []*database.Classroom
		err = json.NewDecoder(resp.Body).Decode(&classrooms)
		assert.NoError(t, err)

		assert.NotNil(t, classrooms)
		assert.NotEmpty(t, classrooms)
	})
}
