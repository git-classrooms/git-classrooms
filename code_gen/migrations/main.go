package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
)

var migraImage = "supabase/migra:3.0.1663481299"
var postgresImage = "postgres:16-alpine"
var postgresPort = "30000"
var postgresUser = "postgres"
var postgresPassword = "postgres"

func main() {
	log.SetOutput(os.Stderr)
	// Start a postgres container...
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("failed to create docker client", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	var eg errgroup.Group
	eg.SetLimit(2)

	eg.Go(func() error {
		return PullImage(ctx, cli, migraImage)
	})
	eg.Go(func() error {
		return PullImage(ctx, cli, postgresImage)
	})
	if err = eg.Wait(); err != nil {
		log.Fatal("failed to pull image", err)
	}

	log.Println("Pulled images")
	log.Println("Starting postgres instance...")

	dbID, err := createPostgresContainer(ctx, cli, postgresImage, postgresPort)
	if err != nil {
		log.Fatal("failed to create postgres container", err)
	}

	var wg sync.WaitGroup
	exitCodeChan := make(chan int, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		if cli.ContainerRemove(context.Background(), dbID, container.RemoveOptions{Force: true, RemoveVolumes: true}); err != nil {
			log.Println("failed to remove container", err)
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		defer cancel()
		defer func() { exitCodeChan <- 1 }()
		// connect to db and create dbs
		db, err := WaitForDBConnection(ctx, cli, dbID)
		if err != nil {
			log.Println("failed to connect to db", err)
			return
		}

		if _, err = db.Exec("CREATE DATABASE gorm;"); err != nil {
			log.Println("failed to create database", err)
			return
		}
		if _, err = db.Exec("CREATE DATABASE goose;"); err != nil {
			log.Println("failed to create database", err)
			return
		}

		var wg errgroup.Group
		wg.SetLimit(2)
		log.Println("Migrating databases")

		// Connect gorm to the postgres container on the first db
		wg.Go(MigrateGormDatabase)
		wg.Go(MigrateGooseDatabase)

		if err := wg.Wait(); err != nil {
			log.Println("failed to migrate databases", err)
			return
		}
		log.Println("migrated databases")

		// docker run -it --rm --network host supabase/migra:3.0.1663481299 migra postgresql://postgres:postgres@localhost:port/goose postgresql://postgres:postgres@localhost:port/postgres
		logOutput, statusCode, err := runMigraContainer(ctx, cli, migraImage)
		if err != nil {
			log.Println("failed to run migra", err)
			return
		}

		if logOutput != nil {
			defer logOutput.Close()
			io.Copy(os.Stdout, logOutput)
		}

		exitCodeChan <- statusCode
	}()

	exitCode := <-exitCodeChan
	wg.Wait()

	cancel()
	os.Exit(exitCode)
}

func MigrateGooseDatabase() error {
	gooseDB, err := GetDBConnection("goose")
	if err != nil {
		return err
	}
	return database.MigrateDatabase(gooseDB)
}

func MigrateGormDatabase() error {
	gormDB, err := gorm.Open(postgres.Open(GetDsn("gorm")), &gorm.Config{})
	if err != nil {
		return err
	}
	db, err := gormDB.DB()
	if err != nil {
		return err
	}
	if err = database.DatabaseStatus(db); err != nil {
		return err
	}
	return MigrateDatabase(gormDB)
}

func MigrateDatabase(db *gorm.DB) error {
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	log.Println("Running database migrations")
	return db.AutoMigrate(
		&database.User{},
		&database.Classroom{},
		&database.Team{},
		&database.UserClassrooms{},
		&database.Assignment{},
		&database.AssignmentProjects{},
		&database.ClassroomInvitation{},
		&database.ManualGradingRubric{},
		&database.ManualGradingResult{},
		&database.AssignmentJunitTest{},
	)
}

func GetDsn(provider string) string {
	return fmt.Sprintf("postgresql://%s:%s@localhost:%s/%s", postgresUser, postgresPassword, postgresPort, provider)
}

func GetDBConnection(provider string) (*sql.DB, error) {
	gorm, err := gorm.Open(postgres.Open(GetDsn(provider)), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db, err := gorm.DB()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func WaitForDBConnection(ctx context.Context, cli *client.Client, containerID string) (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	out, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{ShowStdout: true, Follow: true})
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(out)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "database system is ready to accept connections") {
			break
		}
	}
	out.Close()

	<-time.After(1000 * time.Millisecond)

	db, err := GetDBConnection("postgres")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func PullImage(ctx context.Context, cli *client.Client, imageName string) error {
	log.Println("Pulling image", imageName)
	pullResponse, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return err
	}
	io.Copy(io.Discard, pullResponse)
	defer pullResponse.Close()
	return nil
}

func GetLogs(ctx context.Context, cli *client.Client, contName string) (logOutput io.ReadCloser) {
	options := container.LogsOptions{ShowStdout: true}

	out, err := cli.ContainerLogs(ctx, contName, options)
	if err != nil {
		panic(err)
	}

	return out
}

func createPostgresContainer(ctx context.Context, cli *client.Client, postgresImage string, postgresPort string) (string, error) {
	createResponse, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        postgresImage,
		Env:          []string{"POSTGRES_PASSWORD=postgres"},
		ExposedPorts: nat.PortSet{"5432/tcp": struct{}{}},
		Tty:          true,
	}, &container.HostConfig{PortBindings: nat.PortMap{nat.Port("5432/tcp"): []nat.PortBinding{{HostIP: "127.0.0.1", HostPort: postgresPort}}}}, &network.NetworkingConfig{}, nil, "")
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, createResponse.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	return createResponse.ID, nil
}

func runMigraContainer(ctx context.Context, cli *client.Client, migraImage string) (logOutput io.ReadCloser, statusCode int, err error) {
	statusCode = 1
	createResponse, err := cli.ContainerCreate(ctx, &container.Config{
		Image: migraImage,
		Cmd:   []string{"migra", "--with-privileges", "--unsafe", GetDsn("goose"), GetDsn("gorm")},
		Tty:   true,
	}, &container.HostConfig{NetworkMode: network.NetworkHost, AutoRemove: true}, &network.NetworkingConfig{}, nil, "")
	if err != nil {
		return
	}

	finishChan, errChan := cli.ContainerWait(ctx, createResponse.ID, container.WaitConditionNextExit)

	if err = cli.ContainerStart(ctx, createResponse.ID, container.StartOptions{}); err != nil {
		return
	}

	select {
	case waitReponse := <-finishChan:
		statusCode = int(waitReponse.StatusCode)
		if statusCode == 0 {
			log.Println("No changes detected")
			return
		}
		logOutput = GetLogs(ctx, cli, createResponse.ID)
		log.Println("Migration changes:")
		return

	case err = <-errChan:
		log.Println("failed to run migra", err)
		return
	}
}
