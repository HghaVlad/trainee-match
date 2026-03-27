package postgres

import (
	"context"
	"log/slog"
	"strings"

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

	// don't log polly unsent from outbox
	if sql, ok := data["sql"].(string); ok {
		if strings.Contains(sql, "FROM outbox") &&
			strings.Contains(sql, "sent_at IS NULL") {
			return
		}
	}

	// don't log tx commands cause of polly outbox using it
	if tag, ok := data["commandTag"].(string); ok {
		switch tag {
		case "BEGIN", "COMMIT", "ROLLBACK":
			return
		}
	}

	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}

	slogLevel := slog.LevelDebug
	switch level {
	case tracelog.LogLevelError:
		slogLevel = slog.LevelError
	case tracelog.LogLevelWarn:
		slogLevel = slog.LevelWarn
	case tracelog.LogLevelInfo:
		slogLevel = slog.LevelInfo
	}

	l.logger.LogAttrs(ctx, slogLevel, msg, attrs...)
}
