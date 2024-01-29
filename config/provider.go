package config

import (
	"errors"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var instance *Config
var once sync.Once

func GetConfig() (*Config, error) {
	once.Do(func() {
		path, _ := os.Getwd()

		godotenv.Load(filepath.Join(path, ".env"), filepath.Join(path, ".env.local"))

		instance = &Config{}

		if err := env.Parse(instance); err != nil {
			log.Fatalf("Couldn't parse environment %s", err.Error())
		}
	})

	if instance == nil {
		return nil, errors.New("there is no config instance")
	}
	return instance, nil
}
