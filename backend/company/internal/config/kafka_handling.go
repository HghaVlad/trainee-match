package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

// TODO: maybe each event gets it's own delay and max retry

type KafkaHandling struct {
	BaseRetryDelay time.Duration `env:"KAFKA_CONSUMER_BASE_RETRY_DELAY" envDefault:"50ms" validate:"gt=0"`
	MaxRetries     int           `env:"KAFKA_CONSUMER_MAX_RETRIES"      envDefault:"2"    validate:"gt=0,lte=5"`
}

func LoadKafkaHandlingConfig(validate *validator.Validate) (*KafkaHandling, error) {
	var cfg KafkaHandling

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse kafka config: %w", err)
	}

	if cfg.BaseRetryDelay > maxKafkaConsumerBaseRetryDelay {
		return nil, errors.New("kafka handling base retry must be < 500ms")
	}

	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("validate kafka handling config: %w", err)
	}

	return &cfg, nil
}

const (
	maxKafkaConsumerBaseRetryDelay = time.Millisecond * 500
)
