package database

import (
	"fmt"
	"log/slog"
	"strings"
)

type PsqlConfig struct {
	Host     string `env:"HOST,notEmpty"`
	Port     int    `env:"PORT" envDefault:"5432"`
	Username string `env:"USER,notEmpty"`
	Password string `env:"PASSWORD,notEmpty"`
	Database string `env:"DB,notEmpty"`
}

var _ slog.LogValuer = (*PsqlConfig)(nil)

func (c *PsqlConfig) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("host", c.Host),
		slog.Int("port", c.Port),
		slog.String("username", strings.Repeat("*", 6)),
		slog.String("password", strings.Repeat("*", 6)),
		slog.String("database", c.Database),
	)
}

func (config *PsqlConfig) Dsn() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.Database)
}
