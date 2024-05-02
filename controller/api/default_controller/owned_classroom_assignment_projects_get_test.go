package default_controller

import (
    "context"
    "fmt"
    "net/http/httptest"
    "testing"

    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "gitlab.hs-flensburg.de/gitlab-classroom/model/database"
    "gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
    "gitlab.hs-flensburg.de/gitlab-classroom/utils"
    "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
    "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
    "gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
    gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
    mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestGetOwnedClassroomAssignmentProjects(t *testing.T) {
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
        ID:                 uuid.New(),
        Name:               "Test classroom",
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
        ID:           uuid.New(),
        ClassroomID:  classroom.ID,
        Name:         "Test assignment",
        Description:  "Assignment description",
    }

    err = query.Assignment.WithContext(context.Background()).Create(assignment)
    if err != nil {
        t.Fatalf("could not create test assignment: %s", err.Error())
    }

    project := &database.AssignmentProjects{
        ID:           uuid.New(),
        AssignmentID: assignment.ID,
        UserID:       1,
        ProjectID:    123,
    }

    err = query.AssignmentProjects.WithContext(context.Background()).Create(project)
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
        ctx.SetOwnedClassroom(classroom)
        ctx.SetOwnedClassroomAssignment(assignment)

        fiberContext.Get(c).SetGitlabRepository(gitlabRepo)
        s := session.Get(c)
        s.SetUserState(session.LoggedIn)
        s.SetUserID(1)
        s.Save()
        return c.Next()
    })

    handler := NewApiController(mailRepo)

    t.Run("GetOwnedClassroomAssignmentProjects", func(t *testing.T) {
        app.Get("/api/classrooms/owned/:classroomId/assignments/:assignmentId/projects", handler.GetOwnedClassroomAssignmentProjects)
        route := fmt.Sprintf("/api/classrooms/owned/%s/assignments/%s/projects", classroom.ID.String(), assignment.ID.String())

        req := httptest.NewRequest("GET", route, nil)
        resp, err := app.Test(req)

        assert.Equal(t, fiber.StatusOK, resp.StatusCode)
        assert.NoError(t, err)
        
    })
}
