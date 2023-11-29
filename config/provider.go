package config

type Provider interface {
	GetConfig() *Config
}
