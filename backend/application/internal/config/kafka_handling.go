package config

import "time"

type KafkaHandling struct {
	RetryDelay time.Duration `mapstructure:"retry_delay"`
	RetryCount int           `mapstructure:"retry_count"`
}
