package default_controller

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestOwnedClassroomAssignmentMiddleware(t *testing.T) {
	// ----------------- DB SETUP -------------------
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

	// --------------- END OF DB SETUP ----------------

	classroom := &database.Classroom{
		ID:        uuid.New(),
		Name:      "Test Classroom",
		OwnerID:   1,
		GroupID:   10,
		GroupAccessTokenID: 30,
		GroupAccessToken: "access-token",
	}
	err = query.Classroom.WithContext(context.Background()).Create(classroom)
	if err != nil {
		t.Fatalf("could not create test classroom: %s", err.Error())
	}

	assignment := &database.Assignment{
		ID:          uuid.New(),
		ClassroomID: classroom.ID,
		Name:        "Sample Assignment",
		Description: "A sample assignment for testing.",
	}
	err = query.Assignment.WithContext(context.Background()).Create(assignment)
	if err != nil {
		t.Fatalf("could not create test assignment: %s", err.Error())
	}

	// --------------- END OF SEEDING DATA ----------------
	app := fiber.New()
	session.InitSessionStore(dbURL)

	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetOwnedClassroom(classroom)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewDefaultController()

	t.Run("OwnedClassroomAssignmentMiddleware", func(t *testing.T) {
		app.Get("/api/classrooms/:classroomId/assignments/:assignmentId", handler.OwnedClassroomAssignmentMiddleware, func(c *fiber.Ctx) error {
			return c.SendString("Assignment found")
		})

		route := fmt.Sprintf("/api/classrooms/%s/assignments/%s", classroom.ID, assignment.ID)
		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)
	})
}
