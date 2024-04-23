package default_controller

import (
	"context"
	"encoding/json"
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

func TestGetOwnedClassroomAssignments(t *testing.T) {
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

	classroom := &database.Classroom{
		Name:               "Test classroom",
		OwnerID:            1,
		Description:        "A classroom for testing assignments",
		GroupID:            1,
		GroupAccessTokenID: 20,
		GroupAccessToken:   "token",
	}
	err := query.Classroom.WithContext(context.Background()).Create(classroom)
	if err != nil {
		t.Fatalf("could not create test classroom: %s", err.Error())
	}

	assignments := []*database.Assignment{
		{Title: "Assignment 1", Description: "First test assignment", ClassroomID: classroom.ID},
		{Title: "Assignment 2", Description: "Second test assignment", ClassroomID: classroom.ID},
	}
	for _, assignment := range assignments {
		err = query.Assignment.WithContext(context.Background()).Create(assignment)
		if err != nil {
			t.Fatalf("could not create test assignment: %s", err.Error())
		}
	}

	// ------------ END OF SEEDING DATA -----------------

	session.InitSessionStore(dbURL)
	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetOwnedClassroom(classroom)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("GetOwnedClassroomAssignments", func(t *testing.T) {
		app.Get("/api/classrooms/owned/:classroomId/assignments", handler.GetOwnedClassroomAssignments)
		route := fmt.Sprintf("/api/classrooms/owned/%d/assignments", classroom.ID)

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		retrievedAssignments, err := query.Assignment.WithContext(context.Background()).Where(query.Assignment.ClassroomID.Eq(classroom.ID)).Find()
		assert.NoError(t, err)
		assert.Len(t, retrievedAssignments, len(assignments))
		for i, retrievedAssignment := range retrievedAssignments {
			assert.Equal(t, assignments[i].Title, retrievedAssignment.Title)
			assert.Equal(t, assignments[i].Description, retrievedAssignment.Description)
		}
	})
}
