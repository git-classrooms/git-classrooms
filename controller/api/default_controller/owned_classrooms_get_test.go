package default_controller

import (
	"context"
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

func TestGetOwnedClassrooms(t *testing.T) {
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

	user := &database.User{ID: 1}
	err = query.User.WithContext(context.Background()).Create(user)
	if err != nil {
		t.Fatalf("could not create test user: %s", err.Error())
	}

	classrooms := []*database.Classroom{
		{
			Name:               "Classroom One",
			OwnerID:            1,
			Description:        "Description One",
			GroupID:            10,
			GroupAccessTokenID: 30,
			GroupAccessToken:   "token30",
		},
		{
			Name:               "Classroom Two",
			OwnerID:            1,
			Description:        "Description Two",
			GroupID:            20,
			GroupAccessTokenID: 40,
			GroupAccessToken:   "token40",
		},
	}
	
	for _, classroom := range classrooms {
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

	t.Run("GetOwnedClassrooms", func(t *testing.T) {
		app.Get("/api/classrooms/owned", handler.GetOwnedClassrooms)
		req := httptest.NewRequest("GET", "/api/classrooms/owned", nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		classrooms, err := query.Classroom.WithContext(context.Background()).Where(query.Classroom.OwnerID.Eq(1)).Find()
		assert.NoError(t, err)
		assert.Len(t, classrooms, len(classrooms))
		for i, classroom := range classrooms {
			assert.Equal(t, classrooms[i].Name, classroom.Name)
			assert.Equal(t, classrooms[i].Description, classroom.Description)
			assert.Equal(t, classrooms[i].GroupID, classroom.GroupID)
			assert.Equal(t, classrooms[i].GroupAccessTokenID, classroom.GroupAccessTokenID)
			assert.Equal(t, classrooms[i].GroupAccessToken, classroom.GroupAccessToken)
		}
	})
}