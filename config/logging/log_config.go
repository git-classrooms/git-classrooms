package logging

import (
	"io"
	"log/slog"
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

func (t logType) getHandler(output io.Writer, opts *slog.HandlerOptions) slog.Handler {
	switch t {
	case typeText:
		return slog.NewTextHandler(output, opts)
	case typeJson:
		return slog.NewJSONHandler(output, opts)
	default:
		return slog.NewTextHandler(output, opts)
	}
}

type LogConfig struct {
	Level    logLevel `env:"LEVEL" envDefault:"info"`
	Type     logType  `env:"TYPE" envDefault:"text"`
	Source   bool     `env:"SOURCE" envDefault:"false"`
	FilePath string   `env:"FILE"`
}

var _ slog.LogValuer = (*LogConfig)(nil)

func (c *LogConfig) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("level", string(c.Level)),
		slog.String("type", string(c.Type)),
		slog.Bool("source", c.Source),
	)
}

func (c LogConfig) GetLogger(output io.Writer) *slog.Logger {
	level := c.Level.toSlog()
	opts := &slog.HandlerOptions{
		AddSource: c.Source,
		Level:     level,
	}
	handler := c.Type.getHandler(output, opts)
	return slog.New(handler)
}
