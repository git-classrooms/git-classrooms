package gitlab

import (
	"log/slog"
	"time"
)

type GitlabConfig struct {
	URL          string        `env:"URL"`
	SyncInterval time.Duration `env:"SYNC_INTERVAL" envDefault:"5m"`
}

var _ slog.LogValuer = (*GitlabConfig)(nil)

func (c *GitlabConfig) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("url", c.URL),
		slog.String("syncInterval", c.SyncInterval.String()),
	)
}

func (c *GitlabConfig) GetURL() string {
	return c.URL
}
