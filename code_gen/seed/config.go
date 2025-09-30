package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Postgres PsqlConfig
}

type PsqlConfig struct {
	Hostname string
	Port     int
	Username string
	Password string
	Database string
}

func (c PsqlConfig) Dsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", c.Username, c.Password, c.Hostname, c.Port, c.Database)
}

func ParseConfig() Config {
	godotenv.Load()

	return Config{
		Postgres: PsqlConfig{
			Hostname: "localhost",
			Port:     5432,
			Username: GetEnvDefault("POSTGRES_USER", "postgres"),
			Password: GetEnvDefault("POSTGRES_PASSWORD", "postgres"),
			Database: GetEnvDefault("POSTGRES_DB", "postgres"),
		},
	}
}

func GetEnvDefault(key, def string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return v
}
