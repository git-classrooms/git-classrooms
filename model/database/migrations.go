package database

import (
	"database/sql"
	"embed"
	"log"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func DatabaseStatus(db *sql.DB) error {
	log.Println("Checking database status")

	if err := goose.Status(db, "migrations"); err != nil {
		return (err)
	}
	return nil
}

func MigrateDatabase(db *sql.DB) error {
	log.Println("Running database migrations")

	if err := goose.Up(db, "migrations"); err != nil {
		return (err)
	}

	return nil
}

func init() {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}
}
