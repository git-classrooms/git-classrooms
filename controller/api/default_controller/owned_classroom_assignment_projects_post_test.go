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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestInviteToAssignmentProject(t *testing.T) {
	// Setup the test database with testcontainers
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

	// Seed data
	user := &database.User{ID: 1, GitlabEmail: "user@example.com", Name: "Test User"}
	err = query.User.WithContext(context.Background()).Create(user)
	if err != nil {
		t.Fatalf("could not create test user: %s", err.Error())
	}
	classroom := &database.Classroom{
		Name:               "Test Classroom",
		OwnerID:            1,
		Description:        "Test Description",
		GroupID:            1,
		GroupAccessTokenID: 20,
		GroupAccessToken:   "token",
	}
	err = query.Classroom.WithContext(context.Background()).Create(classroom)
	if err != nil {
		t.Fatalf("could not create test classroom: %s", err.Error())
	}
	assignment := &database.Assignment{
		ClassroomID: classroom.ID,
		Name:        "Test Assignment",
	}
	err = query.Assignment.WithContext(context.Background()).Create(assignment)
	if err != nil {
		t.Fatalf("could not create test assignment: %s", err.Error())
	}

	session.InitSessionStore(dbURL)
	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetOwnedClassroom(classroom)
		ctx.SetOwnedClassroomAssignment(assignment)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("InviteToAssignmentProject", func(t *testing.T) {
		app.Post("/api/classrooms/:classroomId/assignments/:assignmentId/invite", handler.InviteToAssignmentProject)
		route := fmt.Sprintf("/api/classrooms/%s/assignments/%s/invite", classroom.ID.String(), assignment.ID.String())

		req := httptest.NewRequest("POST", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		assert.NoError(t, err)

	})
}
