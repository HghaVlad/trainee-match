package postgres

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/golang-migrate/migrate/v4"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/config"
)

func BuildDBURL(cfg config.Postgres) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		cfg.DBName,
		cfg.SSLMode,
	)
}

func Migrate(cfg config.Postgres) error {
	m, err := migrate.New(
		"file://internal/infra/db/postgres/migrations",
		BuildDBURL(cfg),
	)
	if err != nil {
		return fmt.Errorf("init migrate: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrate: %w", err)
	}

	return nil
}
