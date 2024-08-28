package gitlab

import "time"

type GitlabConfig struct {
	URL          string        `env:"URL"`
	SyncInterval time.Duration `env:"SYNC_INTERVAL" envDefault:"5m"`
}

func (c *GitlabConfig) GetURL() string {
	return c.URL
}
