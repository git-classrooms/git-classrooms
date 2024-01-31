package gitlab

type GitlabConfig struct {
	URL string `env:"URL"`
}

func (c *GitlabConfig) GetURL() string {
	return c.URL
}
