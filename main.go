package main

import (
	"backend/router"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	applicationConfig := config.EnvProvider{}.GetConfig()

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

	router.Routes(app)

	log.Fatal(app.Listen(":3000"))
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
