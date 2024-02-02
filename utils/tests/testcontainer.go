package tests

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartPostgres() (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "postgres",
		},
		WaitingFor: wait.ForAll(wait.ForListeningPort("5432/tcp"), wait.ForLog(".*database system is ready to accept connections.*").AsRegexp().WithStartupTimeout(60+time.Second)).WithStartupTimeoutDefault(60 * time.Second),
	}

	return testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func GetPostgresData(postgres testcontainers.Container, err error) (testcontainers.Container, string, error) {
	if err != nil {
		return nil, "", err
	}

	postgresHost, err := postgres.ContainerIP(context.Background())
	if err != nil {
		return nil, "", err
	}

	return postgres, postgresHost, nil
}
