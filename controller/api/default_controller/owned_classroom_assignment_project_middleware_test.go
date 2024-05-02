package default_controller

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

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

	"github.com/gofiber/fiber/v2"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestOwnedClassroomAssignmentProjectMiddleware(t *testing.T) {
	// Database setup
	t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "false")
	pg, err := tests.StartPostgres()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		pg.Restore(context.Background())
	})
	dbURL, err := pg.ConnectionString(context.Background())

	db, err := gorm.Open(postgresDriver.Open(dbURL), &gorm.Config{})
	if err != nil {
		t.Fatalf("could not connect to database: %s", err.Error())
	}

	utils.MigrateDatabase(db)
	pg.Snapshot(context.Background(), postgres.WithSnapshotName("test-snapshot"))
	query.SetDefault(db)

	// Seeding data
	user := &database.User{ID: 1, GitlabEmail: "user@example.com", Name: "Test User"}
	query.User.WithContext(context.Background()).Create(user)
	classroom := &database.Classroom{
		ID:                 uuid.New(),
		Name:               "Test Classroom",
		OwnerID:            1,
		Description:        "A classroom for testing",
		GroupID:            1,
		GroupAccessTokenID: 20,
		GroupAccessToken:   "token",
	}
	query.Classroom.WithContext(context.Background()).Create(classroom)
	assignment := &database.Assignment{
		ID:          uuid.New(),
		ClassroomID: classroom.ID,
		Name:        "Test Assignment",
		Description: "Assignment for testing",
	}
	query.Assignment.WithContext(context.Background()).Create(assignment)
	assignmentProject := &database.AssignmentProjects{
		ID:           uuid.New(),
		AssignmentID: assignment.ID,
		UserID:       1,
		ProjectID:    100,
	}
	query.AssignmentProjects.WithContext(context.Background()).Create(assignmentProject)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetUserID(1)
		return c.Next()
	})

	app.Use("/api", ctrl.OwnedClassroomAssignmentProjectMiddleware)

	t.Run("ValidAssignmentProjectMiddlewareCall", func(t *testing.T) {
		route := fmt.Sprintf("/api/classrooms/%s/assignments/%s/projects/%s", classroom.ID, assignment.ID, assignmentProject.ID)
		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("InvalidParametersMiddlewareCall", func(t *testing.T) {
		route := fmt.Sprintf("/api/classrooms/%s/assignments/%s/projects/%s", classroom.ID, uuid.New(), uuid.New())
		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)
		assert.Error(t, err)
		assert.NotEqual(t, fiber.StatusOK, resp.StatusCode)
	})
}
