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

func TestJoinedClassroomAssignmentMiddleware(t *testing.T) {
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

	classroomID := uuid.New()
	assignmentID := uuid.New()

	classroom := &database.Classroom{
		ID:                 classroomID,
		Name:               "Test Classroom",
		OwnerID:            1,
		Description:        "Classroom description",
		GroupID:            1,
		GroupAccessTokenID: 20,
		GroupAccessToken:   "token",
	}

	err = query.Classroom.WithContext(context.Background()).Create(classroom)
	if err != nil {
		t.Fatalf("could not create test classroom: %s", err.Error())
	}

	assignment := &database.Assignment{
		ID:                assignmentID,
		ClassroomID:       classroomID,
		TemplateProjectID: 1,
		Name:              "Test Assignment",
		Description:       "Assignment description",
	}

	err = query.Assignment.WithContext(context.Background()).Create(assignment)
	if err != nil {
		t.Fatalf("could not create test assignment: %s", err.Error())
	}

	assignmentProject := &database.AssignmentProjects{
		AssignmentID:       assignmentID,
		UserID:             1,
		AssignmentAccepted: true,
		ProjectID:          1,
	}

	err = query.AssignmentProjects.WithContext(context.Background()).Create(assignmentProject)
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
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("JoinedClassroomAssignmentMiddleware", func(t *testing.T) {
		app.Use("/api/classrooms/:classroomId/assignments/:assignmentId", handler.JoinedClassroomAssignmentMiddleware)
		route := fmt.Sprintf("/api/classrooms/%s/assignments/%s", classroomID.String(), assignmentID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		assignmentProject, err := query.AssignmentProjects.WithContext(context.Background()).
			Where(query.AssignmentProjects.AssignmentID.Eq(assignmentID)).
			First()
		assert.NoError(t, err)
		assert.Equal(t, assignmentID, assignmentProject.AssignmentID)
		assert.Equal(t, 1, assignmentProject.UserID)
	})
}
