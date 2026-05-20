package config

import (
	"errors"
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type ElasticSearch struct {
	Addrs    []string `env:"ELASTIC_NODES" envSeparator:","`
	User     string   `env:"ELASTIC_USER" validate:"required"`
	Password string   `env:"ELASTIC_PASS" validate:"required"`
}

func LoadElasticConfig(validate *validator.Validate) (*ElasticSearch, error) {
	var cfg ElasticSearch

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse elastic config: %w", err)
	}

	if len(cfg.Addrs) == 0 {
		return nil, errors.New("ELASTIC_NODES must not be empty")
	}

	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("validate elastic config: %w", err)
	}

	return &cfg, nil
}
