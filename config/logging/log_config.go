package logging

import (
	"log/slog"
	"os"
)

type logLevel string

const (
	levelDebug logLevel = "debug"
	levelInfo  logLevel = "info"
	levelWarn  logLevel = "warn"
	levelError logLevel = "error"
)

func (l logLevel) toSlog() slog.Level {
	switch l {
	case levelDebug:
		return slog.LevelDebug
	case levelInfo:
		return slog.LevelInfo
	case levelWarn:
		return slog.LevelWarn
	case levelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type logType string

const (
	typeText logType = "text"
	typeJson logType = "json"
)

func (t logType) getHandler(opts *slog.HandlerOptions) slog.Handler {
	switch t {
	case typeText:
		return slog.NewTextHandler(os.Stdout, opts)
	case typeJson:
		return slog.NewJSONHandler(os.Stdout, opts)
	default:
		return slog.NewTextHandler(os.Stdout, opts)
	}
}

type LogConfig struct {
	Level  logLevel `env:"LEVEL" envDefault:"info"`
	Type   logType  `env:"TYPE" envDefault:"text"`
	Source bool     `env:"SOURCE" envDefault:"false"`
}

func (c LogConfig) GetLogger() *slog.Logger {
	level := c.Level.toSlog()
	opts := &slog.HandlerOptions{
		AddSource: c.Source,
		Level:     level,
	}
	handler := c.Type.getHandler(opts)
	return slog.New(handler)
}
