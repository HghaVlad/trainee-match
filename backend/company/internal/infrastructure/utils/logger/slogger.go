package logger

import (
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
