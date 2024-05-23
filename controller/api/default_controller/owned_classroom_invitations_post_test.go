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

func TestInviteToClassroom(t *testing.T) {
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

	classroomQuery := query.Classroom
	testClassRoom := &database.Classroom{
		Name:               "Test classroom",
		OwnerID:            1,
		Description:        "Classroom description",
		GroupID:            1,
		GroupAccessTokenID: 20,
		GroupAccessToken:   "token",
	}

	err = classroomQuery.WithContext(context.Background()).Create(testClassRoom)
	if err != nil {
		t.Fatalf("could not create test classroom: %s", err.Error())
	}

	// ------------ END OF SEEDING DATA -----------------

	session.InitSessionStore(dbURL)
	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetOwnedClassroom(testClassRoom)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("InviteToClassroom", func(t *testing.T) {
		app.Post("/api/classrooms/:classroomId/invitations", handler.InviteToClassroom)
		route := fmt.Sprintf("/api/classrooms/%d/invitations", testClassRoom.ID)

		req := httptest.NewRequest("POST", route, strings.NewReader(`{"memberEmails":["test1@example.com", "test2@example.com"]}`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		invitations, err := query.ClassroomInvitation.WithContext(context.Background()).Where(query.ClassroomInvitation.ClassroomID.Eq(testClassRoom.ID)).Find()
		assert.NoError(t, err)
		assert.Len(t, invitations, 2)
		assert.Equal(t, "test1@example.com", invitations[0].Email)
		assert.Equal(t, "test2@example.com", invitations[1].Email)
	})
}
