package postgres

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
)

func ConnectPgxPoolWithLogger(ctx context.Context, cfg config.Postgres, logger *slog.Logger) (*pgxpool.Pool, error) {
	pgxCfg, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, err
	}

	cc := pgxCfg.ConnConfig

	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		port = 5432
	}

	cc.User = cfg.User
	cc.Password = cfg.Password
	cc.Host = cfg.Host
	cc.Port = uint16(port)
	cc.Database = cfg.Name

	if cfg.SSLMode == "disable" {
		cc.TLSConfig = nil
	} else {
		cc.TLSConfig = &tls.Config{}
	}

	pgxCfg.MaxConns = int32(cfg.MaxPoolConns)
	pgxCfg.MinConns = int32(cfg.MinPoolConns)

	pgxCfg.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   newPgxSlogAdapter(logger),
		LogLevel: tracelog.LogLevelTrace,
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("postgres pgx pool init: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("postgres pgx ping: %w", err)
	}

	return pool, nil
}
