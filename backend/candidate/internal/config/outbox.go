package config

import (
	"github.com/spf13/viper"
	"time"
)

type Outbox struct {
	CandidateTopic string `mapstructure:"CANDIDATE_TOPIC"`
	ResumeTopic    string `mapstructure:"RESUME_TOPIC"`

	MaxAttempts int64 `mapstructure:"MAX_ATTEMPTS"`

	RelayWorkerCount int           `mapstructure:"RELAY_WORKER_COUNT"`
	RelayMinSleep    time.Duration `mapstructure:"RELAY_MIN_SLEEP"`
	RelayMaxSleep    time.Duration `mapstructure:"RELAY_MAX_SLEEP"`
	RelayBatchSize   int           `mapstructure:"RELAY_BATCH_SIZE"`
	BaseRetryDelay   time.Duration `mapstructure:"BASE_RETRY_DELAY"`
	ResetStaleTime   time.Duration `mapstructure:"RESET_STALE_TIME"`
}

// bindOutboxEnv binds outbox related envs to viper keys (kept with Outbox type)
func bindOutboxEnv(v *viper.Viper) error {
	if err := v.BindEnv("OUTBOX.CANDIDATE_TOPIC", "KAFKA_CANDIDATE_TOPIC"); err != nil {
		return err
	}
	if err := v.BindEnv("OUTBOX.RESUME_TOPIC", "KAFKA_RESUME_TOPIC"); err != nil {
		return err
	}
	if err := v.BindEnv("OUTBOX.MAX_ATTEMPTS", "KAFKA_MAX_ATTEMPTS"); err != nil {
		return err
	}
	if err := v.BindEnv("OUTBOX.RELAY_MIN_SLEEP", "OUTBOX_RELAY_MIN_SLEEP"); err != nil {
		return err
	}
	if err := v.BindEnv("OUTBOX.RELAY_MAX_SLEEP", "OUTBOX_RELAY_MAX_SLEEP"); err != nil {
		return err
	}
	if err := v.BindEnv("OUTBOX.RELAY_WORKER_COUNT", "OUTBOX_RELAY_WORKER_COUNT"); err != nil {
		return err
	}
	if err := v.BindEnv("OUTBOX.RELAY_BATCH_SIZE", "OUTBOX_RELAY_BATCH_SIZE"); err != nil {
		return err
	}
	if err := v.BindEnv("OUTBOX.BASE_RETRY_DELAY", "OUTBOX_MESSAGE_BASE_RETRY_DELAY"); err != nil {
		return err
	}
	return nil
}
