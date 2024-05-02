package default_controller

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
	"github.com/google/uuid"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestOwnedClassroomAssignmentMiddleware(t *testing.T) {
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

	// Seeding data
	classroom := &database.Classroom{
		ID: uuid.New(),
		OwnerID: 1,
	}
	err = query.Classroom.WithContext(context.Background()).Create(classroom)
	if err != nil {
		t.Fatalf("could not create test classroom: %s", err.Error())
	}

	assignment := &database.Assignment{
		ID: uuid.New(),
		ClassroomID: classroom.ID,
		Name: "Test Assignment",
	}
	err = query.Assignment.WithContext(context.Background()).Create(assignment)
	if err != nil {
		t.Fatalf("could not create test assignment: %s", err.Error())
	}

	// ------------ END OF SEEDING DATA -----------------

	session.InitSessionStore(dbURL)
	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		fctx := context.Get(c)
		fctx.SetOwnedClassroom(classroom)

		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewDefaultController()

	t.Run("OwnedClassroomAssignmentMiddleware", func(t *testing.T) {
		app.Get("/api/classrooms/:classroomId/assignments/:assignmentId", handler.OwnedClassroomAssignmentMiddleware)
		route := fmt.Sprintf("/api/classrooms/%s/assignments/%s", classroom.ID.String(), assignment.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		queriedAssignment, err := query.Assignment.WithContext(context.Background()).Where(query.Assignment.ID.Eq(assignment.ID)).First()
		assert.NoError(t, err)
		assert.Equal(t, assignment.Name, queriedAssignment.Name)
	})
}
