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
	"golang.org/x/exp/constraints"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func run() error {
	ctx := context.Background()

	gitlabURL := flag.String("gitlabURL", "http://gitlab.localhost:6969", "Base URL of the GitLab instance (default: http://gitlab.localhost:6969)")

	var backupDotenv *bool
	flag.BoolFunc("b", "Backup .env automatically without prompting", func(s string) error {
		if s == "true" {
			backupDotenv = Ptr(true)
		} else {
			backupDotenv = Ptr(false)
		}
		return nil
	})

	flag.Parse()

	config := ParseConfig()
	dsn := config.Postgres.Dsn()

	log.Println("Starting clean infrastructure reset...")

	if err := ExecComposeReset(ctx); err != nil {
		return err
	}

	log.Println("Waiting for PostgreSQL to become available...")

	for {
		if err := ExecPgHealth(ctx); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	log.Println("Start database migration")

	if err := ExecMigrateDB(ctx, dsn); err != nil {
		return err
	}

	log.Println("Start database seeding")

	if err := ExecSeedDB(ctx, dsn); err != nil {
		return err
	}

	gitlabRepo, err := NewGitlabRepo(*gitlabURL, "")
	if err != nil {
		return err
	}

	log.Println("Checking if GitLab is online... This can take a few minutes.")
	current := time.Now()
	for {
		if err := gitlabRepo.Health(ctx); err != nil {
			if time.Since(current) > 30*time.Second {
				current = time.Now()
				log.Printf("GitLab not online yet; retrying… (last error: %v)\n", err)
			}
		} else {
			break
		}
		time.Sleep(5 * time.Second)
	}

	log.Println("Gitlab is online. Creating PAT for root user")
	adminToken, err := gitlabRepo.CreatePersonalAccessTokenForRoot(ctx)
	if err != nil {
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

	if backupDotenv == nil {
		log.Print("Backup .env before updating? (y/n)")
		var char rune
		if _, err := fmt.Scanf("%c", &char); err != nil {
			return err
		}

		backupDotenv = Ptr(char == 'y')
	}

	if *backupDotenv {
		log.Println("Creating .env backup")
		if err := ExecCommand(ctx, "cp .env .env.bak"); err != nil {
			return err
		}
	} else {
		log.Println("Skipping .env backup")
	}

	log.Println("Updating .env with GitLab credentials...")

	if err := UpdateDotenv(map[string]string{
		"AUTH_CLIENT_ID":     application.ApplicationID,
		"AUTH_CLIENT_SECRET": application.Secret,
		"GITLAB_URL":         *gitlabURL,
	}); err != nil {
		return err
	}

	log.Println("Registering GitLab instance runner...")

	gitlabRunner, err := gitlabRepo.CreateInstanceRunner(ctx)
	if err != nil {
		return err
	}

	if err := ExecGitlabRunnerRegister(ctx, *gitlabURL, gitlabRunner.Token); err != nil {
		return err
	}

	log.Println("Instance runner registered")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database connection", err)
	}

	query.SetDefault(db)

	users, err := GetAllUsers(ctx)
	if err != nil {
		return err
	}

	log.Printf("Creating %d users\n", len(users))

	for _, user := range users {
		log.Println("Creating user", user.GitlabEmail)

		gitlabUser, err := gitlabRepo.CreateUser(ctx, user.GitlabUsername, user.GitlabEmail, user.Name)
		if err != nil {
			return err
		}

		if user.ID != gitlabUser.ID {
			return fmt.Errorf("ERROR: IDs are different than expected! DB ID=%d, GitLab ID=%d\n", user.ID, gitlabUser.ID)
		}
	}

	log.Println("All users created")

	owner := users[0]

	log.Println("Retrieving owner access token")

	ownerToken, err := gitlabRepo.GetAccessTokenOfUser(ctx, owner.ID)
	if err != nil {
		return err
	}

	ownerRepo, err := NewGitlabRepo(*gitlabURL, ownerToken.Token)
	if err != nil {
		return err
	}

	log.Printf("Setting up %d template projects\n", len(templateProjects))

	for _, project := range templateProjects {
		log.Println("Creating template project", *project.opts.Name)
		gitlabProject, err := ownerRepo.CreateTemplateProject(ctx, project.opts)
		if err != nil {
			return err
		}

		project.projectID = gitlabProject.ID

		log.Println("Project created with id", gitlabProject.ID)
	}

	classrooms, err := GetAllClassrooms(ctx)
	if err != nil {
		return err
	}

	log.Printf("Creating %d classrooms\n", len(classrooms))

	for _, classroom := range classrooms {
		log.Println("Creating classroom", classroom.Name)

		gitlabClassroom, err := ownerRepo.CreateClassroom(ctx, classroom.Name)
		if err != nil {
			return err
		}

		log.Println("Classroom created with id", gitlabClassroom.ID)

		log.Println("Creating group access token for classroom...")
		groupToken, err := ownerRepo.CreateGroupAccessToken(ctx, gitlabClassroom.ID)
		if err != nil {
			return err
		}

		groupToken2, err := ownerRepo.CreateGroupAccessToken(ctx, gitlabClassroom.ID)
		if err != nil {
			return err
		}

		classroom.GroupID = gitlabClassroom.ID
		classroom.Tokens = []*database.ClassroomToken{
			{
				GroupAccessTokenID:        groupToken.ID,
				GroupAccessToken:          groupToken.Token,
				GroupAccessTokenCreatedAt: *groupToken.CreatedAt,
			}, {
				GroupAccessTokenID:        groupToken2.ID,
				GroupAccessToken:          groupToken2.Token,
				GroupAccessTokenCreatedAt: *groupToken2.CreatedAt,
			},
		}
		if err := SaveClassroom(ctx, classroom); err != nil {
			return err
		}

		if _, err := sqlDB.ExecContext(ctx, "UPDATE classrooms SET group_id = $1 WHERE id = $2", gitlabClassroom.ID, classroom.ID); err != nil {
			return err
		}

		log.Println("Fixing classroom assignment template project IDs...")
		for _, assignment := range classroom.Assignments {
			gitlabProject := templateProjects[abs(assignment.TemplateProjectID)-1]
			if _, err := sqlDB.ExecContext(ctx, "UPDATE assignments SET template_project_id = $1 WHERE id = $2", gitlabProject.projectID, assignment.ID); err != nil {
				return nil
			}
		}

		groupRepo, err := NewGitlabRepo(*gitlabURL, groupToken.Token)
		if err != nil {
			return err
		}

		log.Printf("Adding %d members to classroom\n", len(classroom.Member))

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

		log.Printf("Creating %d teams of classroom\n", len(classroom.Teams))

		for _, team := range classroom.Teams {
			log.Println("Creating team of classroom", team.Name)
			gitlabTeam, err := groupRepo.CreateTeam(ctx, gitlabClassroom.ID, team.Name)
			if err != nil {
				return err
			}

			team.GroupID = gitlabTeam.ID
			if _, err := sqlDB.ExecContext(ctx, "UPDATE teams SET group_id = $1 WHERE id = $2", gitlabTeam.ID, team.ID); err != nil {
				return nil
			}

			log.Printf("Adding %d members to team", len(team.Member))

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

func abs[T constraints.Integer | constraints.Float](v T) T {
	if v < 0 {
		return -v
	}
	return v
}

type templateProject struct {
	opts      *gitlab.CreateProjectOptions
	projectID int
}

var templateProjects = []*templateProject{
	{opts: &gitlab.CreateProjectOptions{
		Name:                 Ptr("Simple Assignment"),
		DefaultBranch:        Ptr("main"),
		InitializeWithReadme: Ptr(true),
		Visibility:           Ptr(gitlab.PublicVisibility),
	}},
	{opts: &gitlab.CreateProjectOptions{
		Name:                 Ptr("Go Assignment"),
		InitializeWithReadme: Ptr(false),
		Visibility:           Ptr(gitlab.PublicVisibility),
	}},
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
