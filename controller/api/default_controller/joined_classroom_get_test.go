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

func TestGetJoinedClassroom(t *testing.T) {
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

	member := &database.User{ID: 2, GitlabEmail: "member@example.com"}
	err = query.User.WithContext(context.Background()).Create(member)

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
		UserID:      member.ID,
		ClassroomID: testClassroom.ID,
		Role:        database.Student,
	}

	err = userClassroomQuery.WithContext(context.Background()).Create(testUserClassroom)
	if err != nil {
		t.Fatalf("could not create user test classroom: %s", err.Error())
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
		s.SetUserID(member.ID)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("GetJoinedClassroom", func(t *testing.T) {
		app.Get("/api/classrooms/joined/:classroomId", handler.GetJoinedClassroom)
		route := fmt.Sprintf("/api/classrooms/joined/%s", testClassroom.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)

		type ClassRoomResponse struct {
			ID          uuid.UUID `json:"id"`
			Name        string    `json:"name"`
			OwnerID     int       `json:"ownerId"`
			Description string    `json:"description"`
			GroupID     int       `json:"groupId"`
		}

		var classRoom *ClassRoomResponse

		err = json.NewDecoder(resp.Body).Decode(&classRoom)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		assert.Equal(t, testClassroom.ID, classRoom.ID)
		assert.Equal(t, testClassroom.Name, classRoom.Name)
		assert.Equal(t, testClassroom.OwnerID, classRoom.OwnerID)
		assert.Equal(t, testClassroom.Description, classRoom.Description)
		assert.Equal(t, testClassroom.GroupID, classRoom.GroupID)
	})
}
