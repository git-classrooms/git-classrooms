package default_controller

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
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
)

func TestJoinedClassroomMiddleware(t *testing.T) {
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

	user := &database.User{ID: 1, GitlabEmail: "test@example.com", Name: "Test User"}
	err = query.User.WithContext(context.Background()).Create(user)
	if err != nil {
		t.Fatalf("could not create test user: %s", err.Error())
	}

	classroom := &database.Classroom{
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

	userClassroom := &database.UserClassrooms{
		UserID:      1,
		ClassroomID: classroom.ID,
		Role:        database.Student,
	}
	err = query.UserClassrooms.WithContext(context.Background()).Create(userClassroom)
	if err != nil {
		t.Fatalf("could not create user classroom association: %s", err.Error())
	}

	// ------------ END OF SEEDING DATA -----------------

	session.InitSessionStore(dbURL)
	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

	app := fiber.New()
	app.Use("/api", func(c *fiber.Ctx) error {
		ctx := fiberContext.Get(c)
		ctx.SetJoinedClassroom(classroom)

		fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
		s := session.Get(c)
		s.SetUserState(session.LoggedIn)
		s.SetUserID(1)
		s.Save()
		return c.Next()
	})

	handler := NewApiController(mailRepo)

	t.Run("JoinedClassroomMiddleware", func(t *testing.T) {
		app.Use("/api/classrooms/joined/:classroomId", handler.JoinedClassroomMiddleware)
		route := fmt.Sprintf("/api/classrooms/joined/%s", classroom.ID.String())

		req := httptest.NewRequest("GET", route, nil)
		resp, err := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NoError(t, err)

		joinedClassroom, err := query.UserClassrooms.WithContext(context.Background()).
			Preload(query.UserClassrooms.Classroom).
			Where(query.UserClassrooms.UserID.Eq(1)).
			First()
		assert.NoError(t, err)
		assert.Equal(t, classroom.Name, joinedClassroom.Classroom.Name)
		assert.Equal(t, classroom.Description, joinedClassroom.Classroom.Description)
		assert.Equal(t, classroom.GroupID, joinedClassroom.Classroom.GroupID)
		assert.Equal(t, classroom.GroupAccessTokenID, joinedClassroom.Classroom.GroupAccessTokenID)
		assert.Equal(t, classroom.GroupAccessToken, joinedClassroom.Classroom.GroupAccessToken)
	})
}
