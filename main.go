//go:generate go run ./code_gen/gorm/main.go
//go:generate swag fmt --exclude frontend
//go:generate swag init --requiredByDefault --exclude frontend
//go:generate mockery
package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/config"
	api "gitlab.hs-flensburg.de/gitlab-classroom/controller/api/default_controller"
	authController "gitlab.hs-flensburg.de/gitlab-classroom/controller/auth"
	"gitlab.hs-flensburg.de/gitlab-classroom/docs"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/httputil"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/mail"
	"gitlab.hs-flensburg.de/gitlab-classroom/router"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/worker"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

var version string = "develop"

//	@title			GitClassrooms – Backend API
//	@version		develop
//	@description	This is the API for our GitClassrooms Webapp.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	GitClassrooms
//	@contact.url	https://git-classrooms.dev
//	@contact.email	info@git-classrooms.dev

//	@license.name	Mozilla Public License 2.0
//	@license.url	https://raw.githubusercontent.com/git-classrooms/git-classrooms/refs/heads/develop/LICENSE

func main() {
	appConfig, err := config.LoadApplicationConfig()
	if err != nil {
		log.Fatal("failed to get application configuration", err)
	}

	setSwaggerInfo(appConfig.PublicURL.String())

	log.Printf("Starting GitClassrooms %s", version)

	mailRepo, err := mail.NewMailRepository(appConfig.PublicURL, appConfig.Mail)
	if err != nil {
		log.Fatal("failed to create mail repository", err)
	}

	db, err := gorm.Open(postgres.Open(appConfig.Database.Dsn()), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get database connection", err)
	}

	session.InitSessionStore(utils.Ptr(appConfig.Database.Dsn()), appConfig.PublicURL)

	if err = database.MigrateDatabase(sqlDB); err != nil {
		log.Fatal("failed to migrate database", err)
	}
	log.Println("DB has been initialized")

	// Set db for gorm-gen
	query.SetDefault(db)

	app := fiber.New(fiber.Config{
		AppName:                 "GitClassrooms",
		ServerHeader:            "GitClassrooms",
		EnableTrustedProxyCheck: len(appConfig.TrustedProxies) > 0,
		TrustedProxies:          appConfig.TrustedProxies,
		ErrorHandler:            errorHandler,
	})

	authCtrl := authController.NewOAuthController(appConfig.Auth, appConfig.GitLab)
	apiController := api.NewApiV1Controller(mailRepo, *appConfig)

	app.Mount("/", router.Routes(authCtrl, apiController, frontendFS, appConfig.Auth))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()
		log.Println("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			log.Println(err)
		}
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.Listen(fmt.Sprintf(":%d", appConfig.Port)); err != nil {
			log.Println(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		dueAssignmentWork := worker.NewDueAssignmentWork(appConfig.GitLab)
		dueAssignmentWorker := worker.NewWorker(dueAssignmentWork)
		dueAssignmentWorker.Start(ctx, 1*time.Minute)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		syncGitlabDbWork := worker.NewSyncGitlabDbWork(appConfig.GitLab, appConfig.PublicURL)
		syncGitlabDbWorker := worker.NewWorker(syncGitlabDbWork)
		syncGitlabDbWorker.Start(ctx, appConfig.GitLab.SyncInterval)
	}()

	wg.Wait()
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	return c.Status(code).JSON(httputil.HTTPError{
		Error: err.Error(),
	})
}

func setSwaggerInfo(appURL string) {
	var schemes []string
	var trimmedAppURL string
	if strings.HasPrefix(appURL, "http://") {
		schemes = []string{"http"}
		trimmedAppURL = strings.TrimPrefix(appURL, "http://")
	} else {
		trimmedAppURL = strings.TrimPrefix(appURL, "https://")
		schemes = []string{"https"}
	}

	docs.SwaggerInfo.Title = "GitClassrooms Backend API"
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Description = "This is the API for our GitClassrooms Webapp."
	docs.SwaggerInfo.Host = trimmedAppURL
	docs.SwaggerInfo.Schemes = schemes
}
