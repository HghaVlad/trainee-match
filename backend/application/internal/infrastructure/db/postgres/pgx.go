package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"

	"github.com/HghaVlad/trainee-match/backend/application/internal/config"
)

func ConnectPgxPoolWithLogger(ctx context.Context, cfg config.DB, logger *slog.Logger) (*pgxpool.Pool, error) {
	pgxCfg, err := buildPgxConfFromAppConf(cfg)
	if err != nil {
		return nil, err
	}

	pgxCfg.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   newPgxSlogAdapter(logger),
		LogLevel: tracelog.LogLevelError,
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

func buildPgxConfFromAppConf(cfg config.DB) (*pgxpool.Config, error) {
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
	cc.Database = cfg.DBName

	if cfg.MaxPoolConns < 0 || cfg.MaxPoolConns > 200 {
		return nil, fmt.Errorf("invalid max pool conns: %d", cfg.MaxPoolConns)
	}

	if cfg.MinPoolConns < 0 || cfg.MinPoolConns > cfg.MaxPoolConns || cfg.MinPoolConns > 200 {
		return nil, fmt.Errorf("invalid min pool conns: %d", cfg.MinPoolConns)
	}

	pgxCfg.MaxConns = int32(cfg.MaxPoolConns)
	pgxCfg.MinConns = int32(cfg.MinPoolConns)

	return pgxCfg, nil
}
