package main

import (
	"database/sql"
	"log"

	// "github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Start a postgres container...

	// connect to db and create dbs
	db, err := sql.Open("postgres", "postgres://postgres:postgres@postgres:30000/postgres")
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	if _, err = db.Exec("CREATE DATABASE gorm;"); err != nil {
		log.Fatal("failed to create database", err)
	}
	if _, err = db.Exec("CREATE DATABASE goose;"); err != nil {
		log.Fatal("failed to create database", err)
	}

	// Connect gorm to the postgres container on the first db
	gormDB, err := gorm.Open(postgres.Open("postgres://postgres:postgres@postgres:30000/gorm"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	err = utils.MigrateDatabase(gormDB)
	if err != nil {
		log.Fatal("failed to migrate database", err)
	}

	// Connect goose to the second postgres container
	gooseDB, err := sql.Open("postgres", "postgres://postgres:postgres@postgres:30000/goose")
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	// setup database

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(gooseDB, "migrations"); err != nil {
		panic(err)
	}

	// docker run -it --rm --network host supabase/migra:3.0.1663481299 migra postgresql://postgres:postgres@localhost:port/goose postgresql://postgres:postgres@localhost:port/postgres
}
