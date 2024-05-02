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
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestGetOwnedClassroomAssignmentProjectGitlab(t *testing.T) {
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

	classroom := &database.Classroom{
		ID:          uuid.New(),
		Name:        "Test Classroom",
		OwnerID:     1,
		Description: "A classroom for testing",
	}
	err = query.Classroom.WithContext(context.Background()).Create(classroom)
	if err != nil {
		t.Fatalf("could not create classroom: %s", err.Error())
	}

	assignment := &database.Assignment{
		ClassroomID: classroom.ID,
		Name:        "Test Assignment",
		Description: "Assignment Description",
	}
	err = query.Assignment.WithContext(context.Background()).Create(assignment)
	if err != nil {
		t.Fatalf("could not create assignment: %s", err.Error())
	}

	project := &database.AssignmentProjects{
		AssignmentID: assignment.ID,
		UserID:       1,
		ProjectID:    12345, // Mock project ID to fetch from GitLab
	}
	err = query.AssignmentProjects.WithContext(context.Background()).Create(project)
	if err != nil {
		t.Fatalf("could not create project: %s", err.Error())
	}

	// ------------ END OF SEEDING DATA -----------------

	session.InitSessionStore(dbURL)
	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetOwnedClassroomAssignmentProject(project)
		ctx.SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("GetOwnedClassroomAssignmentProjectGitlab", func(t *testing.T) {
		app.Get("/api/classrooms/owned/:classroomId/assignments/:assignmentId/gitlab", handler.GetOwnedClassroomAssignmentProjectGitlab)
		route := fmt.Sprintf("/api/classrooms/owned/%s/assignments/%s/gitlab", classroom.ID, assignment.ID)

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusFound, resp.StatusCode)
		assert.NoError(t, err)
	})
}
