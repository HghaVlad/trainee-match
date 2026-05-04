package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type Kafka struct {
	Brokers        []string      `env:"KAFKA_BROKERS"         envSeparator:","`
	ClientID       string        `env:"KAFKA_CLIENT_ID"                        envDefault:"company" validate:"required"`
	ProducerAcks   string        `env:"KAFKA_PRODUCER_ACKS"                    envDefault:"all"     validate:"oneof=none 0 leader 1 all -1"`
	ProducerLinger time.Duration `env:"KAFKA_PRODUCER_LINGER"                  envDefault:"10ms"    validate:"gt=0"`

	ConsumerGroup string `env:"KAFKA_CONSUMER_GROUP" validate:"required"`
	UserTopic     string `env:"KAFKA_USER_TOPIC"     validate:"required"`
}

func LoadKafkaConfig(validate *validator.Validate) (*Kafka, error) {
	var cfg Kafka

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse kafka config: %w", err)
	}

	if len(cfg.Brokers) == 0 {
		return nil, errors.New("KAFKA_BROKERS must not be empty")
	}

	if cfg.ProducerLinger > maxProducerLinger {
		return nil, errors.New("KAFKA_PRODUCER_LINGER should be less than 20ms")
	}

	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("validate kafka config: %w", err)
	}

	return &cfg, nil
}

const (
	maxProducerLinger = 20 * time.Millisecond
)
