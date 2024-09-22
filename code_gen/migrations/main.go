package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var migraImage = "supabase/migra:3.0.1663481299"
var postgresImage = "postgres:16"
var postgresPort = "30000"
var postgresUser = "postgres"
var postgresPassword = "postgres"

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

func PullImage(ctx context.Context, cli *client.Client, imageName string) error {
	log.Println("Pulling image", imageName)
	pullResponse, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return err
	}
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

func main() {
	log.SetOutput(os.Stderr)
	// Start a postgres container...
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("failed to create docker client", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

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

	createResponse, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        postgresImage,
		Env:          []string{"POSTGRES_PASSWORD=postgres"},
		ExposedPorts: nat.PortSet{"5432/tcp": struct{}{}},
	}, &container.HostConfig{PortBindings: nat.PortMap{nat.Port("5432/tcp"): []nat.PortBinding{{HostIP: "127.0.0.1", HostPort: postgresPort}}}}, &network.NetworkingConfig{}, nil, "")
	if err != nil {
		log.Fatal("failed to create container", err)
	}

	if err := cli.ContainerStart(ctx, createResponse.ID, container.StartOptions{}); err != nil {
		log.Fatal("failed to start container", err)
	}

	var wg sync.WaitGroup
	exitCodeChan := make(chan int, 1)
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		if cli.ContainerRemove(context.Background(), createResponse.ID, container.RemoveOptions{Force: true, RemoveVolumes: true}); err != nil {
			log.Println("failed to remove container", err)
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		db := func() *sql.DB {
			ticker := time.NewTicker(1 * time.Second)
			for {
				select {
				case <-ctx.Done():
					return nil
				case <-ticker.C:
					db, err := GetDBConnection("postgres")
					if err != nil {
						log.Println("trying to connect to postgres", err)
						continue
					}
					return db
				}
			}
		}()

		// connect to db and create dbs
		defer cancel()

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

		wg.Go(func() error {
			// Connect gorm to the postgres container on the first db
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
			return utils.MigrateDatabase(gormDB)
		})

		wg.Go(func() error {
			// Connect goose to the second postgres container
			gooseDB, err := GetDBConnection("goose")
			if err != nil {
				return err
			}
			return database.MigrateDatabase(gooseDB)
		})

		if err := wg.Wait(); err != nil {
			log.Println("failed to migrate databases", err)
			return
		}
		log.Println("migrated databases")

		// docker run -it --rm --network host supabase/migra:3.0.1663481299 migra postgresql://postgres:postgres@localhost:port/goose postgresql://postgres:postgres@localhost:port/postgres
		createResponse, err := cli.ContainerCreate(ctx, &container.Config{
			Image: migraImage,
			Cmd:   []string{"migra", "--with-privileges", "--unsafe", GetDsn("goose"), GetDsn("gorm")},
			Tty:   true,
		}, &container.HostConfig{NetworkMode: network.NetworkHost, AutoRemove: true}, &network.NetworkingConfig{}, nil, "")
		if err != nil {
			log.Println("failed to create container", err)
			return
		}

		finishChan, errChan := cli.ContainerWait(ctx, createResponse.ID, container.WaitConditionNextExit)

		if err := cli.ContainerStart(ctx, createResponse.ID, container.StartOptions{}); err != nil {
			log.Println("failed to start container", err)
			return
		}

		select {
		case waitReponse := <-finishChan:
			exitCodeChan <- int(waitReponse.StatusCode)
			if waitReponse.StatusCode == 0 {
				log.Println("No changes detected")
				return
			}
			logOutput := GetLogs(ctx, cli, createResponse.ID)
			log.Println("Migration changes:")

			io.Copy(os.Stdout, logOutput)
		case err := <-errChan:
			log.Println("failed to run migra", err)
		}
	}()

	exitCode := <-exitCodeChan
	wg.Wait()
	os.Exit(exitCode)
}
