//go:generate go run ./code_gen/gorm/main.go
//go:generate swag fmt --exclude frontend,controller/api
//go:generate swag init --requiredByDefault --exclude frontend,controller/api
//go:generate mockery
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/config"
	apiController "gitlab.hs-flensburg.de/gitlab-classroom/controller/api/default_controller"
	api "gitlab.hs-flensburg.de/gitlab-classroom/controller/api_v2/default_controller"
	authController "gitlab.hs-flensburg.de/gitlab-classroom/controller/auth"
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

var version string = "development"

//	@title			GitLab Classrooms – Backend API
//	@version		1.0.0
//	@description	This is the API for our Gitlab Classroom Webapp
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Hochschule Flensburg – Applied Computer Science (Master)
//	@contact.url	https://hs-flensburg.de

//	@license.name	Proprietary
//	@license.url	https://gitlab.hs-flensburg.de/fb3-masterprojekt-gitlab-classroom/gitlab-classroom/-/raw/develop/LICENSE.md

func main() {
	appConfig, err := config.LoadApplicationConfig()
	if err != nil {
		log.Fatal("failed to get application configuration", err)
	}

	log.Println("SentryEnabled", appConfig.Sentry.IsEnabled())
	if appConfig.Sentry.IsEnabled() {
		err = sentry.Init(sentry.ClientOptions{
			Dsn:         appConfig.Sentry.GetDSN(),
			Environment: appConfig.Sentry.GetEnv(),
			Release:     version,

			// Enable printing of SDK debug messages.
			// Useful when getting started or trying to figure something out.
			Debug: true,
		})

		if err != nil {
			log.Fatalf("failed to init sentry: %s", err)
		}

		defer sentry.Flush(2 * time.Second)
	}

	mailRepo, err := mail.NewMailRepository(appConfig.PublicURL, appConfig.Mail)
	if err != nil {
		log.Fatal("failed to create mail repository", err)
	}

	db, err := gorm.Open(postgres.Open(appConfig.Database.Dsn()), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	session.InitSessionStore(utils.Ptr(appConfig.Database.Dsn()))

	err = utils.MigrateDatabase(db)
	if err != nil {
		log.Fatal("failed to migrate database", err)
	}
	log.Println("DB has been initialized")

	// Set db for gorm-gen
	query.SetDefault(db)

	app := fiber.New(fiber.Config{
		EnableTrustedProxyCheck: true,
		TrustedProxies:          appConfig.TrustedProxies,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			return c.Status(code).JSON(httputil.HTTPError{
				Error: err.Error(),
			})
		},
	})

	authCtrl := authController.NewOAuthController(appConfig.Auth, appConfig.GitLab)
	apiCtrl := apiController.NewApiController(mailRepo)
	v2Controller := api.NewApiV2Controller(mailRepo, *appConfig)

	router.Routes(app, authCtrl, apiCtrl, v2Controller, appConfig.FrontendPath, appConfig.Auth)

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

	wg.Wait()
}
