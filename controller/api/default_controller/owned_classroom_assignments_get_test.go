package default_controller

import (
	"context"
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

	owner := &database.User{ID: 1, GitlabEmail: "owner@example.com"}
	err = query.User.WithContext(context.Background()).Create(owner)
	if err != nil {
		t.Fatalf("could not create test user: %s", err.Error())
	}

	testClassroom := &database.Classroom{
		Name:               "Test classroom",
		OwnerID:            owner.ID,
		Description:        "Classroom description",
		GroupID:            1,
		GroupAccessTokenID: 20,
		GroupAccessToken:   "token",
	}

	err = query.Classroom.WithContext(context.Background()).Create(testClassroom)
	if err != nil {
		t.Fatalf("could not create test classroom: %s", err.Error())
	}

	testClassroomAssignments := []*database.Assignment{
		{
			ClassroomID:       testClassroom.ID,
			TemplateProjectID: 1,
			Name:              "Test Assignment 1",
			Description:       "Description 1",
		},
		{
			ClassroomID:       testClassroom.ID,
			TemplateProjectID: 2,
			Name:              "Test Assignment 2",
			Description:       "Description 2",
		},
	}

	for _, a := range testClassroomAssignments {
		err = query.Assignment.WithContext(context.Background()).Create(a)
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
		ctx.SetOwnedClassroom(testClassroom)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(owner.ID)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("GetOwnedClassroomAssignments", func(t *testing.T) {
		app.Get("/classrooms/owned/:classroomId/assignments", handler.GetOwnedClassroomAssignments)
		route := fmt.Sprintf("/api/classrooms/owned/%s/assignments", testClassroom.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)

		var classroomAssignments []*database.Assignment
		err = json.NewDecoder(resp.Body).Decode(&classroomAssignments)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Len(t, classroomAssignments, len(testClassroomAssignments))

		for i, assignment := range classroomAssignments {
			assert.Equal(t, testClassroomAssignments[i].ID, assignment.ID)
			assert.Equal(t, testClassroomAssignments[i].ClassroomID, assignment.ClassroomID)
			assert.Equal(t, testClassroomAssignments[i].TemplateProjectID, assignment.TemplateProjectID)
			assert.Equal(t, testClassroomAssignments[i].Name, assignment.Name)
			assert.Equal(t, testClassroomAssignments[i].Description, assignment.Description)
		}
	})
}