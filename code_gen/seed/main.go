package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"
)

func run() error {
	ctx := context.Background()

	gitlabURL := flag.String("gitlabURL", "http://gitlab.localhost:6969", "URL of the gitlab-instance")
	adminToken := flag.String("adminToken", "", "Access Token of the instance Admin")
	notStart := flag.Bool("notStart", false, "Access Token of the instance Admin")

	flag.Parse()

	if *notStart && *adminToken == "" {
		return fmt.Errorf("admin token has to be provided")
	}

	if !*notStart {
		if err := ExecComposeReset(ctx); err != nil {
			return err
		}

		gitlabRepo, err := NewGitlabRepo(*gitlabURL, "")
		if err != nil {
			return err
		}

		for {
			online, err := gitlabRepo.Health(ctx)
			if err != nil {
				log.Println("Not yet online", err)
			} else if !online {
				log.Println("Not yet online")
			} else {
				break
			}
			time.Sleep(5 * time.Second)
		}

		log.Println("Open the following link in your browser, login and create an accessToken with the following permissions:")
		log.Println("write_repository, api, create_runner, admin_mode, sudo")
		url := fmt.Sprintf("%s/-/user_settings/personal_access_tokens?page=1&state=active&sort=expires_asc", *gitlabURL)
		log.Println(url)
		ExecCommand(ctx, fmt.Sprintf("open %s", url))

		fmt.Print("Admin Token: ")
		if _, err := fmt.Scanln(adminToken); err != nil {
			return err
		}
	}

	gitlabRepo, err := NewGitlabRepo(*gitlabURL, *adminToken)
	if err != nil {
		return err
	}

	log.Println("Creating application for GitClassrooms")

	application, err := gitlabRepo.CreateApplication(ctx)
	if err != nil {
		return err
	}

	log.Println("Updating dotenv")

	if err := UpdateDotenv(map[string]string{
		"AUTH_CLIENT_ID":     application.ApplicationID,
		"AUTH_CLIENT_SECRET": application.Secret,
	}); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
