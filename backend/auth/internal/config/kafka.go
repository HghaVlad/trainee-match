package config

import "time"

type Kafka struct {
	Brokers        []string      `mapstructure:"BROKERS"`
	ClientID       string        `mapstructure:"CLIENT_ID"`
	ProducerAcks   string        `mapstructure:"PRODUCER_ACKS"   validate:"oneof=none 0 leader 1 all -1"`
	ProducerLinger time.Duration `mapstructure:"PRODUCER_LINGER" validate:"gt=0"`
	UserTopic      string        `mapstructure:"USER_TOPIC"`
}
