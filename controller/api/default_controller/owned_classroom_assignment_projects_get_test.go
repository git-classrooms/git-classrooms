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

func TestGetOwnedClassroomAssignmentProjects(t *testing.T) {
	// DB setup and migration
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

	// Seeding data for user, classroom, and assignment
	user := &database.User{ID: 1}
	err = query.User.WithContext(context.Background()).Create(user)
	if err != nil {
		t.Fatalf("could not create test user: %s", err.Error())
	}

	testClassRoom := &database.Classroom{
		Name:    "Test classroom",
		OwnerID: 1,
	}

	err = query.Classroom.WithContext(context.Background()).Create(testClassRoom)
	if err != nil {
		t.Fatalf("could not create test classroom: %s", err.Error())
	}

	testAssignment := &database.Assignment{
		ClassroomID: testClassRoom.ID,
		Name:        "Sample Assignment",
	}

	err = query.Assignment.WithContext(context.Background()).Create(testAssignment)
	if err != nil {
		t.Fatalf("could not create test assignment: %s", err.Error())
	}

	testProject := &database.AssignmentProjects{
		AssignmentID: testAssignment.ID,
		UserID:       1,
	}

	err = query.AssignmentProjects.WithContext(context.Background()).Create(testProject)
	if err != nil {
		t.Fatalf("could not create test project: %s", err.Error())
	}

	// Session and context setup
	session.InitSessionStore(dbURL)
	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetOwnedClassroom(testClassRoom)
		ctx.SetOwnedClassroomAssignment(testAssignment)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("GetOwnedClassroomAssignmentProjects", func(t *testing.T) {
		app.Get("/api/classrooms/owned/:classroomId/assignments/:assignmentId/projects", handler.GetOwnedClassroomAssignmentProjects)
		route := fmt.Sprintf("/api/classrooms/owned/%s/assignments/%s/projects", testClassRoom.ID, testAssignment.ID)

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

	})
}
