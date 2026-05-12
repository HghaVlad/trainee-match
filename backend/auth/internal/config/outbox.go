package config

import "time"

type Outbox struct {
	UserTopic        string        `mapstructure:"USER_TOPIC"`
	BaseRetryDelay   time.Duration `mapstructure:"BASE_RETRY_DELAY"`
	MaxRetries       int           `mapstructure:"MAX_RETRIES"`
	BatchSize        int           `mapstructure:"BATCH_SIZE"`
	RelayMinSleep    time.Duration `mapstructure:"RELAY_MIN_SLEEP"`
	RelayMaxSleep    time.Duration `mapstructure:"RELAY_MAX_SLEEP"`
	RelayWorkerCount int           `mapstructure:"RELAY_WORKER_COUNT"`
}
