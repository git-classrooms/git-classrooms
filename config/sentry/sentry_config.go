package sentry

type SentryConfig struct {
	Dsn     string `env:"DSN"`
	Enabled bool   `env:"ENABLED"`
	Env     string `env:"ENV"`
}

func (c *SentryConfig) GetDSN() string {
	return c.Dsn
}

func (c *SentryConfig) IsEnabled() bool {
	return c.Enabled
}

func (c *SentryConfig) GetEnv() string{
	return c.Env
}

