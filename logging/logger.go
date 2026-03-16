package logging

import (
	"log/slog"
	"os"
	"strings"
)

func GetDefaultLogger() *slog.Logger {
	envType := os.Getenv("ENV")
	var handler slog.Handler
	if strings.ToLower(envType) == "prod" {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		handler = slog.NewTextHandler(os.Stdout, nil)
	}
	defaultLogger := slog.New(handler)
	return defaultLogger
}
