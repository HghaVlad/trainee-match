package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	HTTP           HTTP
	Postgres       Postgres
	Redis          Redis
	Kafka          Kafka
	Outbox         Outbox
	KafkaHandling KafkaHandling
	SchemaRegistry SchemaRegistry
}

func Load() (*Config, error) {
	validate := validator.New()

	httpCfg, err := LoadHTTPConfig(validate)
	if err != nil {
		return nil, err
	}

	postgresCfg, err := LoadPostgresConfig(validate)
	if err != nil {
		return nil, err
	}

	redisCfg, err := LoadRedisConfig(validate)
	if err != nil {
		return nil, err
	}

	kafkaCfg, err := LoadKafkaConfig(validate)
	if err != nil {
		return nil, err
	}

	outboxCfg, err := LoadOutboxConfig(validate)
	if err != nil {
		return nil, err
	}

	kafkaHandling, err := LoadKafkaHandlingConfig(validate)
	if err != nil {
		return nil, err
	}

	schemaRegCfg, err := LoadSchemaRegistryConfig(validate)
	if err != nil {
		return nil, err
	}

	return &Config{
		HTTP:           *httpCfg,
		Postgres:       *postgresCfg,
		Redis:          *redisCfg,
		SchemaRegistry: *schemaRegCfg,
		Kafka:          *kafkaCfg,
		KafkaHandling:  *kafkaHandling,
		Outbox:         *outboxCfg,
	}, nil
}

func loadConfigSection[T any](validate *validator.Validate, sectionName string) (*T, error) {
	var cfg T

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse %s config: %w", sectionName, err)
	}

	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("validate %s config: %w", sectionName, err)
	}

	return &cfg, nil
}
