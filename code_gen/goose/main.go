package main

import (
	"context"
	"log"
	"os"

	"github.com/pressly/goose/v3"
	_ "gorm.io/driver/postgres"

	"gitlab.hs-flensburg.de/gitlab-classroom/config"
)

func main() {
	args := os.Args

	cfg, err := config.LoadApplicationConfig()
	if err != nil {
		log.Printf("failed to load application config: %v", err)
		return
	}

	command := args[1]
	db, err := goose.OpenDBWithDriver("postgres", cfg.Database.Dsn())

	if err != nil {
		log.Printf("goose: failed to open DB: %v\n", err)
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("goose: failed to close DB: %v\n", err)
		}
	}()

	arguments := []string{}
	if len(args) > 2 {
		arguments = append(arguments, args[2:]...)
	}

	if err := goose.RunContext(context.Background(), command, db, "model/database/migrations", arguments...); err != nil {
		log.Printf("goose %v: %v", command, err)
		return
	}
}
