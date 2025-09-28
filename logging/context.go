package logging

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type contextKey string

const logKey contextKey = "logger"

func SetFiberContextLogger(ctx *fiber.Ctx, logger *slog.Logger) {
	ctx.Locals(logKey, logger)
}

func GetFiberContextLogger(ctx *fiber.Ctx) *slog.Logger {
	return ctx.Locals(logKey).(*slog.Logger)
}

func SetLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, logKey, logger)
}

func GetLogger(ctx context.Context) *slog.Logger {
	return ctx.Value(logKey).(*slog.Logger)
}
