package main

import (
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"

	"backend/model/database"
	"backend/model/database/query"
)

type PsqlConfig struct {
	Host     string `env:"HOST,notEmpty"`
	Port     int    `env:"PORT,notEmpty" env_default:"5432"`
	Username string `env:"USER,notEmpty"`
	Password string `env:"PASSWORD,notEmpty"`
	Database string `env:"DB,notEmpty"`
}

func (config PsqlConfig) Dsn() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.Database)
}

type ApplicationConfig struct {
	Database PsqlConfig `envPrefix:"POSTGRES_"`
}

func main() {
	if _, err := os.Stat(".env.local"); err == nil {
		godotenv.Load(".env", ".env.local")
	} else {
		godotenv.Load()
	}

	config := ApplicationConfig{}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("Couldn't parse environment %s", err.Error())
	}

	db, err := gorm.Open(postgres.Open(config.Database.Dsn()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	log.Println("Running database migrations")
	err = db.AutoMigrate(
		&database.User{},
		&database.Classroom{},
		&database.UserClassrooms{},
		&database.Assignment{},
		&database.AssignmentProjects{},
	)

	// Uncomment this to generate Model Code if the Model changed
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

	log.Fatal(app.Listen(":3000"))
}

func generateGormGen(db *gorm.DB) {
	g := gen.NewGenerator(gen.Config{
		OutPath: "model/database/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(db)

	g.ApplyBasic(&database.User{},
		&database.Classroom{},
		&database.UserClassrooms{},
		&database.Assignment{},
		&database.AssignmentProjects{},
	)

	g.Execute()
}
