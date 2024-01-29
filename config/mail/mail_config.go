package mail

type Config struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
}
