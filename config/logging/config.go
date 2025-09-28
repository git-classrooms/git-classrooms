package logging

import "log/slog"

type Config interface {
	GetLogger() (*slog.Logger, error)
}
