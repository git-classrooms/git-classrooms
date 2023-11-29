package config

import "backend/config/database"

type Config struct {
	Database database.PsqlConfig `envPrefix:"POSTGRES_"`
}
