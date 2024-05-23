package default_controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
)

type testJoinedClassroomAssignmentResponse struct {
	AssignmentProjects *database.AssignmentProjects `json:"assignmentProjects"`
	ProjectPath        string                       `json:"projectPath"`
}

func testJoinedClassroomAssignmentQuery(classroomID uuid.UUID, c *fiber.Ctx) *gorm.DB {
	// Mock implementation of the query function.
	db := query.GetDB()
	return db.Where("classroom_id = ?", classroomID)
}

func TestGetJoinedClassroomAssignments(t *testing.T) {
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

	classroomQuery := query.Classroom
	testClassRoom := &database.Classroom{
		Name:               "Test classroom",
		OwnerID:            1,
		Description:        "Classroom description",
		GroupID:            1,
		GroupAccessTokenID: 20,
		GroupAccessToken:   "token",
	}

	err = classroomQuery.WithContext(context.Background()).Create(testClassRoom)
	if err != nil {
		t.Fatalf("could not create test classroom: %s", err.Error())
	}

	assignmentQuery := query.Assignment
	testAssignment := &database.Assignment{
		ClassroomID:       testClassRoom.ID,
		TemplateProjectID: 1,
		Name:              "Test Assignment",
		Description:       "Assignment description",
	}

	err = assignmentQuery.WithContext(context.Background()).Create(testAssignment)
	if err != nil {
		t.Fatalf("could not create test assignment: %s", err.Error())
	}

	testAssignmentProject := &database.AssignmentProjects{
		AssignmentID:       testAssignment.ID,
		UserID:             1,
		AssignmentAccepted: true,
		ProjectID:          1,
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
	app.Use("/api", func(c *fiber.Cctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetJoinedClassroom(testClassRoom)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("GetJoinedClassroomAssignments", func(t *testing.T) {
		app.Get("/api/classrooms/joined/:classroomId/assignments", handler.GetJoinedClassroomAssignments)
		route := fmt.Sprintf("/api/classrooms/joined/%s/assignments", testClassRoom.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		var responses []*testJoinedClassroomAssignmentResponse
		err = json.NewDecoder(resp.Body).Decode(&responses)
		assert.NoError(t, err)

		assert.Len(t, responses, 1)
		assert.Equal(t, testAssignmentProject.ID, responses[0].AssignmentProjects.ID)
		assert.Equal(t, "/api/v1/classrooms/owned/1/gitlab", responses[0].ProjectPath)
	})
}
