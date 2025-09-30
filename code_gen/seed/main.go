package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/xanzy/go-gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func run() error {
	ctx := context.Background()

	gitlabURL := flag.String("gitlabURL", "http://gitlab.localhost:6969", "URL of the gitlab-instance")

	flag.Parse()

	config := ParseConfig()
	dsn := config.Postgres.Dsn()

	if err := ExecComposeReset(ctx); err != nil {
		return err
	}

	for {
		if err := ExecPgHealth(ctx); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err := ExecMigrateDB(ctx, dsn); err != nil {
		return err
	}

	if err := ExecSeedDB(ctx, dsn); err != nil {
		return err
	}

	gitlabRepo, err := NewGitlabRepo(*gitlabURL, "")
	if err != nil {
		return err
	}

	log.Println("Checking if gitlab is online. This can take a few minutes.")
	current := time.Now()
	for {
		if err := gitlabRepo.Health(ctx); err != nil {
			if time.Since(current) > 30*time.Second {
				current = time.Now()
				log.Println("Still not online waiting", err)
			}
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

	var adminToken string
	fmt.Print("Admin Token: ")
	if _, err := fmt.Scanln(&adminToken); err != nil {
		return err
	}

	gitlabRepo, err = NewGitlabRepo(*gitlabURL, adminToken)
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
		"GITLAB_URL":         *gitlabURL,
	}); err != nil {
		return err
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get database connection", err)
	}

	query.SetDefault(db)

	users, err := GetAllUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		log.Println("Creating user", user.GitlabEmail)

		gitlabUser, err := gitlabRepo.CreateUser(ctx, user.GitlabUsername, user.GitlabEmail, user.Name)
		if err != nil {
			return err
		}

		if user.ID != gitlabUser.ID {
			return fmt.Errorf("ERROR: IDs are different then expected! DB-ID: %d, GitlabID: %d\n", user.ID, gitlabUser.ID)
		}
	}

	log.Println("All users created")

	classrooms, err := GetAllClassrooms(ctx)
	if err != nil {
		return err
	}

	for _, classroom := range classrooms {
		log.Println("Creating classroom", classroom.Name)

		ownerToken, err := gitlabRepo.GetAccessTokenOfUser(ctx, classroom.OwnerID)
		if err != nil {
			return err
		}

		ownerRepo, err := NewGitlabRepo(*gitlabURL, ownerToken.Token)
		if err != nil {
			return err
		}

		gitlabClassroom, err := ownerRepo.CreateClassroom(ctx, classroom.Name)
		if err != nil {
			return err
		}

		log.Println("classroom created with id", gitlabClassroom.ID)

		log.Println("Getting accessToken for classroom")
		groupToken, err := ownerRepo.CreateGroupAccessToken(ctx, gitlabClassroom.ID)
		if err != nil {
			return err
		}

		classroom.GroupID = gitlabClassroom.ID
		classroom.GroupAccessTokenID = groupToken.ID
		classroom.GroupAccessToken = groupToken.Token
		classroom.GroupAccessTokenCreatedAt = *groupToken.CreatedAt
		if err := SaveClassroom(ctx, classroom); err != nil {
			return err
		}

		if _, err := sqlDB.ExecContext(ctx, "UPDATE classrooms SET group_id = $1 WHERE id = $2", gitlabClassroom.ID, classroom.ID); err != nil {
			return err
		}

		groupRepo, err := NewGitlabRepo(*gitlabURL, groupToken.Token)
		if err != nil {
			return err
		}

		for _, member := range classroom.Member {
			log.Println("Adding member to classroom", member.User.GitlabEmail)

			if member.UserID == classroom.OwnerID {
				continue
			}

			permission := gitlab.GuestPermissions
			if classroom.StudentsViewAllProjects || member.Role == database.Moderator {
				permission = gitlab.ReporterPermissions
			} else if member.Role == database.Owner {
				permission = gitlab.OwnerPermissions
			}
			if err := groupRepo.AddUserToGroup(ctx, gitlabClassroom.ID, member.UserID, permission); err != nil {
				return err
			}
		}

		for _, team := range classroom.Teams {
			log.Println("creating team of classroom", team.Name)
			gitlabTeam, err := groupRepo.CreateTeam(ctx, gitlabClassroom.ID, team.Name)
			if err != nil {
				return err
			}

			team.GroupID = gitlabTeam.ID
			if _, err := sqlDB.ExecContext(ctx, "UPDATE teams SET group_id = $1 WHERE id = $2", gitlabTeam.ID, team.ID); err != nil {
				return nil
			}

			for _, member := range team.Member {
				log.Println("Adding member to team", member.User.GitlabEmail)

				if err := groupRepo.AddUserToGroup(ctx, gitlabTeam.ID, member.UserID, gitlab.ReporterPermissions); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
