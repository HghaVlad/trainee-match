package postgres

import (
	"context"
	"fmt"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, conf *config.DB) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.DbName)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
