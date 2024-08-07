package gitlab

type GitlabConfig struct {
	URL                 string `env:"URL"`
	SyncIntervalSeconds int    `env:"SYNC_INTERVAL_SECONDS" envDefault:"300"`
}

func (c *GitlabConfig) GetURL() string {
	return c.URL
}
