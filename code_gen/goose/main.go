package main

import (
	"context"

	"github.com/pressly/goose/v3"
	"gitlab.hs-flensburg.de/gitlab-classroom/config"
	_ "gorm.io/driver/postgres"

	"log"
	"os"
)

func main() {
	args := os.Args

	cfg, err := config.LoadApplicationConfig()
	if err != nil {
		log.Fatalf("failed to load application config: %v", err)
	}

	command := args[1]
	db, err := goose.OpenDBWithDriver("postgres", cfg.Database.Dsn())
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	arguments := []string{}
	if len(args) > 2 {
		arguments = append(arguments, args[2:]...)
	}

	if err := goose.RunContext(context.Background(), command, db, "model/database/migrations", arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
