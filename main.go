package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	apiController "gitlab.hs-flensburg.de/gitlab-classroom/controller/api"
	authController "gitlab.hs-flensburg.de/gitlab-classroom/controller/auth"
	"gitlab.hs-flensburg.de/gitlab-classroom/router"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
	"log"

	"gitlab.hs-flensburg.de/gitlab-classroom/config"

	dbModel "gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
)

func main() {
	appConfig, err := config.LoadApplicationConfig()
	if err != nil {
		log.Fatal("failed to get application configuration", err.Error())
	}

	db, err := gorm.Open(postgres.Open(appConfig.Database.Dsn()), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	log.Println("Running database migrations")
	err = db.AutoMigrate(
		&dbModel.User{},
		&dbModel.Classroom{},
		&dbModel.UserClassrooms{},
		&dbModel.Assignment{},
		&dbModel.AssignmentProjects{},
	)

	// Uncomment this to generate Query Code if the Model changed
	// generateGormGen(db)

	if err != nil {
		panic("failed to migrate database")
	}
	log.Println("DB has been initialized")

	// Set db for gorm-gen
	query.SetDefault(db)

	app := fiber.New()

	authCtrl := authController.NewOAuthController(appConfig.Auth, appConfig.GitLab)
	apiCtrl := apiController.NewApiController()

	router.Routes(app, authCtrl, apiCtrl, appConfig.FrontendPath, appConfig.Auth)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", appConfig.Port)))
}

//lint:ignore U1000 Ignore unused function to generate Query Code if the Model changed
func generateGormGen(db *gorm.DB) {
	g := gen.NewGenerator(gen.Config{
		OutPath: "model/database/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(db)

	g.ApplyBasic(
		&dbModel.User{},
		&dbModel.Classroom{},
		&dbModel.UserClassrooms{},
		&dbModel.Assignment{},
		&dbModel.AssignmentProjects{},
	)

	g.Execute()
}
