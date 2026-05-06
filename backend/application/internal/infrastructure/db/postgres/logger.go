package postgres

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/tracelog"
)

type pgxSlogAdapter struct {
	logger *slog.Logger
}

func newPgxSlogAdapter(l *slog.Logger) *pgxSlogAdapter {
	return &pgxSlogAdapter{logger: l}
}

func (l *pgxSlogAdapter) Log(
	ctx context.Context,
	level tracelog.LogLevel,
	msg string,
	data map[string]any,
) {
	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}

	slogLevel := slog.LevelDebug
	switch level {
	case tracelog.LogLevelNone:
		return
	case tracelog.LogLevelError:
		slogLevel = slog.LevelError
	case tracelog.LogLevelWarn:
		slogLevel = slog.LevelWarn
	case tracelog.LogLevelInfo:
		slogLevel = slog.LevelInfo
	case tracelog.LogLevelDebug, tracelog.LogLevelTrace:
		slogLevel = slog.LevelDebug
	}

	l.logger.LogAttrs(ctx, slogLevel, msg, attrs...)
}
