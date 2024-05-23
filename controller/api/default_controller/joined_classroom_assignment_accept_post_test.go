package default_controller

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	testcontainerspostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestJoinAssignment(t *testing.T) {
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

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		t.Fatalf("could not connect to database: %s", err.Error())
	}

	err = utils.MigrateDatabase(db)
	if err != nil {
		t.Fatalf("could not migrate database: %s", err.Error())
	}

	// Take a snapshot of the database
	err = pg.Snapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	query.SetDefault(db)

	// ------------ END OF DB SETUP -----------------

	// Create test user
	user := &database.User{ID: 1, GitlabEmail: "test@example.com", Name: "Test User"}
	err = query.User.WithContext(context.Background()).Create(user)
	if err != nil {
		t.Fatalf("could not create test user: %s", err.Error())
	}

	// Create test classroom
	classroomQuery := query.Classroom
	classroomID := uuid.New()
	testClassroom := &database.Classroom{
		ID:                 classroomID,
		Name:               "Test classroom",
		OwnerID:            1,
		Description:        "Classroom description",
		GroupID:            1,
		GroupAccessTokenID: 20,
		GroupAccessToken:   "token",
	}
	err = classroomQuery.WithContext(context.Background()).Create(testClassroom)
	if err != nil {
		t.Fatalf("could not create test classroom: %s", err.Error())
	}

	// Create test assignment
	assignmentID := uuid.New()
	testAssignment := &database.Assignment{
		ID:                assignmentID,
		ClassroomID:       classroomID,
		TemplateProjectID: 1234,
		Name:              "Test Assignment",
		Description:       "Test Assignment Description",
	}
	err = query.Assignment.WithContext(context.Background()).Create(testAssignment)
	if err != nil {
		t.Fatalf("could not create test assignment: %s", err.Error())
	}

	// Create test assignment project
	testAssignmentProject := &database.AssignmentProjects{
		AssignmentID:       assignmentID,
		UserID:             1,
		AssignmentAccepted: false,
	}
	err = query.AssignmentProjects.WithContext(context.Background()).Create(testAssignmentProject)
	if err != nil {
		t.Fatalf("could not create test assignment project: %s", err.Error())
	}

	// ------------ END OF SEEDING DATA -----------------

	session.InitSessionStore(dbURL)
	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetJoinedClassroom(&database.UserClassrooms{ClassroomID: classroomID})
		ctx.SetUserID(1)
		ctx.SetGitlabRepository(gitlabRepo)

		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("JoinAssignment", func(t *testing.T) {
		app.Post("/api/classrooms/joined/:classroomId/assignments/:assignmentId/accept", handler.JoinAssignment)
		route := fmt.Sprintf("/api/classrooms/joined/%s/assignments/%s/accept", classroomID.String(), assignmentID.String())

		req := httptest.NewRequest("POST", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusAccepted, resp.StatusCode)
		assert.NoError(t, err)

		assignmentProject, err := query.AssignmentProjects.WithContext(context.Background()).Where(query.AssignmentProjects.AssignmentID.Eq(assignmentID)).Where(query.AssignmentProjects.UserID.Eq(1)).First()
		assert.NoError(t, err)
		assert.True(t, assignmentProject.AssignmentAccepted)
	})
}
