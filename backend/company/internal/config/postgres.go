package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Postgres struct {
	Host     string `env:"POSTGRES_HOST"     validate:"required"`
	Port     string `env:"POSTGRES_PORT"     validate:"required,numeric"                                         envDefault:"5432"`
	Name     string `env:"POSTGRES_DB"       validate:"required"`
	User     string `env:"POSTGRES_USER"     validate:"required"`
	Password string `env:"POSTGRES_PASSWORD" validate:"required"`
	SSLMode  string `env:"POSTGRES_SSL_MODE" validate:"oneof=disable allow prefer require verify-ca verify-full" envDefault:"disable"`

	MaxPoolConns int `env:"POSTGRES_MAX_POOL_CONNS" envDefault:"10" validate:"gte=1,lte=100"`
	MinPoolConns int `env:"POSTGRES_MIN_POOL_CONNS" envDefault:"2"  validate:"gte=0,lte=100"`
}

func LoadPostgresConfig(validate *validator.Validate) (*Postgres, error) {
	cfg, err := loadConfigSection[Postgres](validate, "postgres")
	if err != nil {
		return nil, err
	}

	if cfg.MinPoolConns > cfg.MaxPoolConns {
		return nil, fmt.Errorf(
			"validate postgres config: POSTGRES_MIN_POOL_CONNS (%d) can't be greater than POSTGRES_MAX_POOL_CONNS (%d)",
			cfg.MinPoolConns,
			cfg.MaxPoolConns,
		)
	}

	return cfg, nil
}
