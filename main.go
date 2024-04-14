//go:generate go run ./code_gen/gorm/main.go
//go:generate swag fmt
//go:generate swag init
//go:generate mockery
package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/config"
	apiController "gitlab.hs-flensburg.de/gitlab-classroom/controller/api/default_controller"
	authController "gitlab.hs-flensburg.de/gitlab-classroom/controller/auth"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/httputil"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/mail"
	"gitlab.hs-flensburg.de/gitlab-classroom/router"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var version string = "development"

//	@title			Gitlab Classroom API
//	@version		1.0
//	@description	This is the API for our Gitlab Classroom Webapp
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	GPL 3.0 | MIT | Apache 2.0 | 3-BSD
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/api/v1

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
				Error:   err.Error(),
				Success: false,
			})
		},
	})

	authCtrl := authController.NewOAuthController(appConfig.Auth, appConfig.GitLab)
	apiCtrl := apiController.NewApiController(mailRepo)

	router.Routes(app, authCtrl, apiCtrl, appConfig.FrontendPath, appConfig.Auth)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", appConfig.Port)))
}
