//go:generate go run ./code_gen/gorm_gen.go
//go:generate mockery
package main

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/config"
	"gitlab.hs-flensburg.de/gitlab-classroom/context/session"
	apiController "gitlab.hs-flensburg.de/gitlab-classroom/controller/api/default_controller"
	authController "gitlab.hs-flensburg.de/gitlab-classroom/controller/auth"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/mail"
	"gitlab.hs-flensburg.de/gitlab-classroom/router"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	appConfig, err := config.LoadApplicationConfig()
	if err != nil {
		log.Fatal("failed to get application configuration", err)
	}

	mailRepo, err := mail.NewMailRepository(appConfig.PublicURL, appConfig.Mail)
	if err != nil {
		log.Fatal("failed to create mail repository", err)
	}

	db, err := gorm.Open(postgres.Open(appConfig.Database.Dsn()), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	session.InitSessionStore(appConfig.Database.Dsn())

	err = utils.MigrateDatabase(db)
	if err != nil {
		log.Fatal("failed to migrate database", err)
	}
	log.Println("DB has been initialized")

	// Set db for gorm-gen
	query.SetDefault(db)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error":   err.Error(),
				"success": false,
			})
		},
	})

	authCtrl := authController.NewOAuthController(appConfig.Auth, appConfig.GitLab)
	apiCtrl := apiController.NewApiController(mailRepo)

	router.Routes(app, authCtrl, apiCtrl, appConfig.FrontendPath, appConfig.Auth)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", appConfig.Port)))
}
