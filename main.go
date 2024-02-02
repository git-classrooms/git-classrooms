//go:generate go run ./code_gen/gorm_gen.go
//go:generate mockery
package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/config"
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

	err = utils.MigrateDatabase(db)
	if err != nil {
		log.Fatal("failed to migrate database", err)
	}
	log.Println("DB has been initialized")

	// Set db for gorm-gen
	query.SetDefault(db)

	app := fiber.New()

	authCtrl := authController.NewOAuthController(appConfig.Auth, appConfig.GitLab)
	apiCtrl := apiController.NewApiController(mailRepo)

	router.Routes(app, authCtrl, apiCtrl, appConfig.FrontendPath, appConfig.Auth)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", appConfig.Port)))
}
