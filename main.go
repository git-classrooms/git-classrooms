//go:generate go run ./code_gen/gorm_gen.go
//go:generate mockery
package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/config"
	apiController "gitlab.hs-flensburg.de/gitlab-classroom/controller/api/default_controller"
	authController "gitlab.hs-flensburg.de/gitlab-classroom/controller/auth"
	dbModel "gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/mail"
	"gitlab.hs-flensburg.de/gitlab-classroom/router"
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

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	log.Println("Running database migrations")
	err = db.AutoMigrate(
		&dbModel.User{},
		&dbModel.Classroom{},
		&dbModel.UserClassrooms{},
		&dbModel.Assignment{},
		&dbModel.AssignmentProjects{},
		&dbModel.ClassroomInvitation{},
	)
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
