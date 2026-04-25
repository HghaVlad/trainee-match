package logger

import (
	"context"
	"log/slog"
	"os"
)

func NewSlogLogger() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger
}

type loggerCtxKeyT struct{}

//nolint:gochecknoglobals // ctx key
var loggerCtxKey = loggerCtxKeyT{}

func FromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerCtxKey).(*slog.Logger)
	if !ok {
		panic("logger not found in context. logger middleware is not applied")
	}

	return logger
}

func WithLoggerContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}
