package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type EnvProvider struct {
	instance *Config
	once     sync.Once
}

func (p EnvProvider) GetConfig() *Config {
	p.once.Do(func() {
		path, _ := os.Getwd()

		godotenv.Load(filepath.Join(path, ".env"), filepath.Join(path, ".env.local"))

		p.instance = &Config{}

		if err := env.Parse(p.instance); err != nil {
			log.Fatalf("Couldn't parse environment %s", err.Error())
		}
	})
	return p.instance
}
