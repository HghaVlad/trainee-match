package config

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

type Outbox struct {
	VacancyTopic       string `env:"KAFKA_VACANCY_TOPIC"        validate:"required"`
	CompanyMemberTopic string `env:"KAFKA_COMPANY_MEMBER_TOPIC" validate:"required"`
	CompanyTopic       string `env:"KAFKA_COMPANY_TOPIC"        validate:"required"`

	BaseRetryDelay   time.Duration `env:"OUTBOX_MESSAGE_BASE_RETRY_DELAY" envDefault:"5s"   validate:"gt=0"`
	MaxRetries       int           `env:"OUTBOX_MAX_RETRIES"              envDefault:"5"    validate:"gt=0,lte=8"`
	BatchSize        int           `env:"OUTBOX_MESSAGE_BATCH_SIZE"       envDefault:"100"  validate:"gte=5,lte=400"`
	RelayMinSleep    time.Duration `env:"OUTBOX_RELAY_MIN_SLEEP"          envDefault:"50ms"`
	RelayMaxSleep    time.Duration `env:"OUTBOX_RELAY_MAX_SLEEP"          envDefault:"5s"`
	RelayWorkerCount int           `env:"OUTBOX_RELAY_WORKER_COUNT"       envDefault:"3"    validate:"gte=1,lte=8"`
}

func LoadOutboxConfig(validate *validator.Validate) (*Outbox, error) {
	cfg, err := loadConfigSection[Outbox](validate, "outbox")
	if err != nil {
		return nil, err
	}

	if cfg.BaseRetryDelay < minBaseRetryDelay || cfg.BaseRetryDelay > maxBaseRetryDelay {
		return nil, errors.New("validate outbox config: MESSAGE_BASE_RETRY_DELAY must be between 4 and 10s")
	}

	if cfg.RelayMinSleep < minRelaySleep || cfg.RelayMinSleep > maxRelaySleep {
		return nil, errors.New("validate outbox config: OUTBOX_RELAY_SLEEP must be between 10 and 500ms")
	}

	if cfg.RelayMaxSleep < minMaxRelaySleep || cfg.RelayMaxSleep > maxMaxRelaySleep {
		return nil, errors.New("validate outbox config: OUTBOX_RELAY_MAX_SLEEP must be between 1 and 10s")
	}

	return cfg, nil
}

const (
	minBaseRetryDelay = 4 * time.Second
	maxBaseRetryDelay = 10 * time.Second
)

const (
	minRelaySleep = 5 * time.Millisecond
	maxRelaySleep = 500 * time.Millisecond
)

const (
	minMaxRelaySleep = 1 * time.Second
	maxMaxRelaySleep = 10 * time.Second
)
