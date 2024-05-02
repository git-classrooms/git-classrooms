package default_controller

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetOwnedClassroomInvitations(t *testing.T) {
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
		ID: uuid.New(),
		Name: "Test classroom",
		OwnerID: 1,
		Description: "Classroom description",
		GroupID: 1,
		GroupAccessTokenID: 20,
		GroupAccessToken: "token",
	}

	err = classroomQuery.WithContext(context.Background()).Create(testClassRoom)
	if err != nil {
		t.Fatalf("could not create test classroom: %s", err.Error())
	}

	invitation := &database.ClassroomInvitation{
		ID: uuid.New(),
		ClassroomID: testClassRoom.ID,
		Email: "test@example.com",
		ExpiryDate: time.Now().Add(24 * time.Hour),
		Status: database.ClassroomInvitationPending,
	}

	err = query.ClassroomInvitation.WithContext(context.Background()).Create(invitation)
	if err != nil {
		t.Fatalf("could not create test invitation: %s", err.Error())
	}

	// ------------ END OF SEEDING DATA -----------------

	session.InitSessionStore(dbURL)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetOwnedClassroom(testClassRoom)

		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(nil)

	t.Run("GetOwnedClassroomInvitations", func(t *testing.T) {
		app.Get("/api/classrooms/owned/:classroomId/invitations", handler.GetOwnedClassroomInvitations)
		route := fmt.Sprintf("/api/classrooms/owned/%s/invitations", testClassRoom.ID)

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		invitations, err := query.ClassroomInvitation.WithContext(context.Background()).Where(query.ClassroomInvitation.ClassroomID.Eq(testClassRoom.ID)).Find()
		assert.NoError(t, err)
		assert.NotEmpty(t, invitations)
		assert.Equal(t, "test@example.com", invitations[0].Email)
	})
}
