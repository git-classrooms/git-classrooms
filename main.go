package main

import (
	"backend/context"
	"backend/handler"
	"backend/router"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
	"log"

	"backend/config"

	dbModel "backend/model/database"
	"backend/model/database/query"
)

func main() {
	applicationConfig, err := config.GetConfig()
	if err != nil {
		panic("failed to get application configuration")
	}

	db, err := gorm.Open(postgres.Open(applicationConfig.Database.Dsn()), &gorm.Config{})

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

	app.Get("/api/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	router.Routes(app, applicationConfig)

	app.Use("/api", handler.AuthMiddleware)
	app.Get("/api/secret", func(c *fiber.Ctx) error {
		repo := context.GetGitlabRepository(c)
		user, err := repo.GetCurrentUser()
		if err != nil {
			return err
		}

		return c.JSON(user)
	})

	app.Get("/api/*", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNotFound) })

	// we need to redirect all other routes to the frontend
	spaFile := fmt.Sprintf("%s/index.html", applicationConfig.FrontendPath)
	app.Get("*", func(c *fiber.Ctx) error { return c.SendFile(spaFile) })

	log.Fatal(app.Listen(fmt.Sprintf(":%d", applicationConfig.Port)))
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
