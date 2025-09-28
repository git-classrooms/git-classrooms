package mail

import (
	"log/slog"
	"strings"
)

type MailConfig struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
}

var _ slog.LogValuer = (*MailConfig)(nil)

func (c *MailConfig) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("host", c.Host),
		slog.Int("port", c.Port),
		slog.String("username", strings.Repeat("*", 6)),
		slog.String("password", strings.Repeat("*", 6)),
	)
}

func (c *MailConfig) GetHost() string {
	return c.Host
}

func (c *MailConfig) GetPort() int {
	return c.Port
}

func (c *MailConfig) GetUser() string {
	return c.User
}

func (c *MailConfig) GetPassword() string {
	return c.Password
}
