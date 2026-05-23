package config

import (
	"time"

	"github.com/spf13/viper"
)

type Kafka struct {
	Brokers  []string `mapstructure:"BROKERS"`
	ClientId string   `mapstructure:"CLIENT_ID"`

	ProducerAcks   string        `mapstructure:"PRODUCER_ACKS"`
	ProducerLinger time.Duration `mapstructure:"PRODUCER_LINGER"`
}

// bindKafkaEnv binds kafka related envs to viper keys (kept with kafka types)
func bindKafkaEnv(v *viper.Viper) error {
	if err := v.BindEnv("KAFKA.BROKERS", "KAFKA_BROKERS"); err != nil {
		return err
	}
	if err := v.BindEnv("KAFKA.CLIENT_ID", "KAFKA_CLIENT_ID"); err != nil {
		return err
	}
	if err := v.BindEnv("KAFKA.PRODUCER_ACKS", "KAFKA_PRODUCER_ACKS"); err != nil {
		return err
	}
	if err := v.BindEnv("KAFKA.PRODUCER_LINGER", "KAFKA_PRODUCER_LINGER"); err != nil {
		return err
	}
	return nil
}
