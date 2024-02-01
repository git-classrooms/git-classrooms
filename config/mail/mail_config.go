package mail

type MailConfig struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
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
