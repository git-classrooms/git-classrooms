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
	apiControllerMock "gitlab.hs-flensburg.de/gitlab-classroom/controller/api/_mock"
	authController "gitlab.hs-flensburg.de/gitlab-classroom/controller/auth"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepo "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
	mailRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/mail/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/router"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	db_tests "gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IntegrationTest struct {
	dbURL        string
	container    *postgres.PostgresContainer
	snapshotName string
}

var integrationTest IntegrationTest

func TestMain(m *testing.M) {
	integrationTest = IntegrationTest{
		snapshotName: "integration-test",
	}

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

	// 1. Migrate database
	err = utils.MigrateDatabase(db)
	if err != nil {
		log.Fatalf("could not migrate database: %s", err.Error())
	}

	// close the database connection to create the snapshot
	sqlDB, _ := db.DB()
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
	session.InitSessionStore(nil)

	integrationTest.container = pg
	integrationTest.dbURL = dbURL

	code := m.Run()

	pg.Terminate(context.Background())
	os.Exit(code)
}

func newJsonRequest(route string, object any, httpType string) *http.Request {
	jsonData, err := json.Marshal(object)
	if err != nil {
		log.Fatalf("could not create json of object: %s", object)
	}

	req := httptest.NewRequest(httpType, route, bytes.NewReader(jsonData))
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

func setupApp(t *testing.T, user *database.User, gitlabRepo gitlabRepo.Repository) *fiber.App {
	mailRepo := mailRepoMock.NewMockRepository(t)
	session.InitSessionStore(&integrationTest.dbURL)

	app := fiber.New()

	apiCtrl := apiControllerMock.NewMockController(t)
	v2Controller := NewApiV2Controller(mailRepo, config.ApplicationConfig{})
	authCtrl := authController.NewTestAuthController(user, gitlabRepo)

	redirectUrl, _ := url.Parse("http://example.com")

	router.Routes(app, authCtrl, apiCtrl, v2Controller, "public", &auth.OAuthConfig{RedirectURL: redirectUrl})

	return app
}
