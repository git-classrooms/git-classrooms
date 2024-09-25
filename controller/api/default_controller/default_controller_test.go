package api

import (
	"context"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gitlab.hs-flensburg.de/gitlab-classroom/config"
	"gitlab.hs-flensburg.de/gitlab-classroom/config/auth"
	authController "gitlab.hs-flensburg.de/gitlab-classroom/controller/auth"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/router"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
)

const testURL = "http://example.com"

type IntegrationTest struct {
	dbURL        string
	container    *postgres.PostgresContainer
	snapshotName string
	publicURL    *url.URL
}

var integrationTest IntegrationTest

func TestMain(m *testing.M) {
	integrationTest = IntegrationTest{
		snapshotName: "integration-test",
	}

	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	pg, err := db_tests.StartPostgres()
	if err != nil {
		log.Fatalf("Failed to start postgres container: %s", err.Error())
	}

	dbURL, err := pg.ConnectionString(context.Background())
	if err != nil {
		log.Fatalf("Failed to obtain connection string: %s", err.Error())
	}

	log.Printf("DBURL: %s", dbURL)

	db, err := gorm.Open(postgresDriver.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("could not connect to database: %s", err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("could not get database connection: %s", err.Error())
	}

	// 1. Migrate database
	err = database.MigrateDatabase(sqlDB)
	if err != nil {
		log.Fatalf("could not migrate database: %s", err.Error())
	}

	// close the database connection to create the snapshot
	sqlDB.Close()

	// 2. Create a snapshot of the database to restore later
	err = pg.Snapshot(context.Background(), postgres.WithSnapshotName(integrationTest.snapshotName))
	if err != nil {
		log.Fatalf("Could not create database snapshot: %s", err.Error())
	}

	// open the database connection again
	db, err = gorm.Open(postgresDriver.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("could not connect to database: %s", err.Error())
	}

	query.SetDefault(db)
	publicURL, _ := url.Parse(testURL)
	session.InitSessionStore(nil, publicURL)

	integrationTest.publicURL = publicURL
	integrationTest.container = pg
	integrationTest.dbURL = dbURL

	code := m.Run()

	pg.Terminate(context.Background())
	os.Exit(code)
}

func restoreDatabase(t *testing.T) {
	t.Log("Restore database snapshot...")
	err := integrationTest.container.Restore(context.Background())
	if err != nil {
		t.Fatalf("could not restore container snapshot: %s", err.Error())
	}
}

func setupApp(t *testing.T, user *database.User) (*fiber.App, *gitlabRepoMock.MockRepository, *mailRepoMock.MockRepository) {
	gitlabRepo := gitlabRepoMock.NewMockRepository(t)
	mailRepo := mailRepoMock.NewMockRepository(t)
	session.InitSessionStore(nil, integrationTest.publicURL)

	session.CsrfConfig.Next = func(c *fiber.Ctx) bool { return true }

	app := fiber.New()

	apiController := NewAPIV1Controller(mailRepo, config.ApplicationConfig{PublicURL: integrationTest.publicURL})
	authCtrl := authController.NewTestAuthController(user, gitlabRepo)

	router.Routes(app, authCtrl, apiController, "public", &auth.OAuthConfig{RedirectURL: integrationTest.publicURL})

	return app, gitlabRepo, mailRepo
}

func saveClassroom(t *testing.T, classroom *database.Classroom) {
	err := query.Classroom.WithContext(context.Background()).Save(classroom)
	if err != nil {
		t.Fatal("Could not save classroom!")
	}
}
