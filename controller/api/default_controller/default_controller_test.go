package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

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
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const testUrl = "http://example.com"

type IntegrationTest struct {
	dbURL        string
	container    *postgres.PostgresContainer
	snapshotName string
	publicUrl    *url.URL
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

	// open the database connection agian
	db, err = gorm.Open(postgresDriver.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("could not connect to database: %s", err.Error())
	}

	query.SetDefault(db)
	publicUrl, _ := url.Parse(testUrl)
	session.InitSessionStore(nil, publicUrl)

	integrationTest.publicUrl = publicUrl
	integrationTest.container = pg
	integrationTest.dbURL = dbURL

	code := m.Run()

	pg.Terminate(context.Background())
	os.Exit(code)
}

func newJsonRequest(route string, object any, method string) *http.Request {
	jsonData, err := json.Marshal(object)
	if err != nil {
		log.Fatalf("could not create json of object: %s", object)
	}

	req := httptest.NewRequest(method, route, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	return req
}

func newPostJsonRequest(route string, object any) *http.Request {
	return newJsonRequest(route, object, "POST")
}

func newPutJsonRequest(route string, object any) *http.Request {
	return newJsonRequest(route, object, "PUT")
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
	session.InitSessionStore(nil, integrationTest.publicUrl)

	session.CsrfConfig.Next = func(c *fiber.Ctx) bool { return true }

	apiController := NewApiV1Controller(mailRepo, config.ApplicationConfig{PublicURL: integrationTest.publicUrl})
	authCtrl := authController.NewTestAuthController(user, gitlabRepo)

	app := router.Routes(authCtrl, apiController, os.DirFS("."), &auth.OAuthConfig{RedirectURL: integrationTest.publicUrl})

	return app, gitlabRepo, mailRepo
}

func saveClassroom(t *testing.T, classroom *database.Classroom) {
	err := query.Classroom.WithContext(context.Background()).Save(classroom)
	if err != nil {
		t.Fatal("Could not save classroom!")
	}
}
