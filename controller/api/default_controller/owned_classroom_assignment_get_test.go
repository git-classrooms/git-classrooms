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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestGetOwnedClassroomAssignment(t *testing.T) {
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

	classroomQuery := query.Classroom
	testClassroom := &database.Classroom{
		Name:               "Test classroom",
		OwnerID:            owner.ID,
		Description:        "Classroom description",
		GroupID:            1,
		GroupAccessTokenID: 20,
		GroupAccessToken:   "token",
	}

	err = classroomQuery.WithContext(context.Background()).Create(testClassroom)
	if err != nil {
		t.Fatalf("could not create test classroom: %s", err.Error())
	}

	userClassroomQuery := query.UserClassrooms
	testUserClassroom := &database.UserClassrooms{
		UserID:      owner.ID,
		ClassroomID: testClassRoom.ID,
		Role:        database.Student,
	}

	err = userClassroomQuery.WithContext(context.Background()).Create(testUserClassroom)
	if err != nil {
		t.Fatalf("could not create user test classroom: %s", err.Error())
	}

	classroomAssignmentQuery := query.Assignment
	testClassroomAssignment := &database.Assignment{
		ClassroomID:       testClassroom.ID,
		TemplateProjectID: 1,
		Name:              "Test classroom assignment",
		Description:       "Classroom assignment description",
	}

	err = classroomAssignmentQuery.WithContext(context.Background()).Create(testClassroomAssignment)
	if err != nil {
		t.Fatalf("could not create test classroom assignment: %s", err.Error())
	}

	// ------------ END OF SEEDING DATA -----------------

	session.InitSessionStore(dbURL)
	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetOwnedClassroom(testClassRoom)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(owner.ID)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("GetOwnedClassroomAssignment", func(t *testing.T) {
		app.Get("classrooms/owned/:classroomId/assignments/:assignmentId", handler.GetOwnedClassroomAssignment)
		route := fmt.Sprintf("/api/classrooms/owned/%d/assignments/%d", testClassRoom.ID, testClassroomAssignment.ID)

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)

		type ClassroomAssignmentResponse struct {
			ID                uuid.UUID  `json:"id"`
			CreatedAt         time.Time  `json:"createdAt"`
			UpdatedAt         time.Time  `json:"updatedAt"`
			ClassroomID       uuid.UUID  `json:"classroomId"`
			TemplateProjectID int        `json:"templateProjectId"`
			Name              string     `json:"name"`
			Description       string     `json:"description"`
			DueDate           *time.Time `json:"dueDate"`
		}

		var classroomAssignment *ClassroomAssignmentResponse

		err = json.NewDecoder(resp.Body).Decode(&classroomAssignment)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		assert.Equal(t, testClassroomAssignment.ID, classroomAssignment.ID)
		assert.Equal(t, testClassroomAssignment.ClassroomID, classroomAssignment.ClassroomID)
		assert.Equal(t, testClassroomAssignment.TemplateProjectID, classroomAssignment.TemplateProjectID)
		assert.Equal(t, testClassroomAssignment.Name, classroomAssignment.Name)
		assert.Equal(t, testClassroomAssignment.Description, classroomAssignment.Description)
	})
}
