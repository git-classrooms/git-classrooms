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
    fiberContext "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
    "gitlab.hs-flensburg.de/gitlab-classroom/utils"
    "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
    "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
    postgresDriver "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestOwnedClassroomMiddleware(t *testing.T) {
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

    user := &database.User{ID: 1, Name: "Test User", GitlabEmail: "test@example.com"}
    err = query.User.WithContext(context.Background()).Create(user)
    if err != nil {
        t.Fatalf("could not create test user: %s", err.Error())
    }

    testClassRoom := &database.Classroom{
        Name:               "Test classroom",
        OwnerID:            1,
        Description:        "Classroom description",
        GroupID:            1,
        GroupAccessTokenID: 20,
        GroupAccessToken:   "token",
    }

    err = query.Classroom.WithContext(context.Background()).Create(testClassRoom)
    if err != nil {
        t.Fatalf("could not create test classroom: %s", err.Error())
    }

    // ------------ END OF SEEDING DATA -----------------

    session.InitSessionStore(dbURL)
	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)

    app := fiber.New()
    app.Use("/api", func(c *fiber.Ctx) error {
        ctx := fiberContext.Get(c)
        ctx.SetUserID(1)
        ctx.SetSession(session.Get(c))
        return c.Next()
    })

    handler := NewApiController(mailRepo)

    t.Run("OwnedClassroomMiddleware", func(t *testing.T) {
        app.Get("/api/classrooms/owned/:classroomId", handler.OwnedClassroomMiddleware, func(c *fiber.Ctx) error {
            classroom := fiberContext.Get(c).GetOwnedClassroom()
            return c.JSON(classroom)
        })

        route := fmt.Sprintf("/api/classrooms/owned/%s", testClassRoom.ID.String())

        req := httptest.NewRequest("GET", route, nil)
        resp, err := app.Test(req)

        assert.Equal(t, fiber.StatusOK, resp.StatusCode)
        assert.NoError(t, err)

        retrievedClassroom, err := query.Classroom.WithContext(context.Background()).Where(query.Classroom.OwnerID.Eq(1)).First()
        assert.NoError(t, err)
        assert.Equal(t, testClassRoom.Name, retrievedClassroom.Name)
    })
}
