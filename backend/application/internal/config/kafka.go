package config

import "time"

type Kafka struct {
	Brokers  []string `mapstructure:"BROKERS"`
	ClientID string   `mapstructure:"CLIENT_ID"`

	ConsumerGroup  string   `mapstructure:"CONSUMER_GROUP"`
	ConsumerTopics []string `mapstructure:"CONSUMER_TOPICS"`

	ProducerAcks   string        `mapstructure:"PRODUCER_ACKS"`
	ProducerLinger time.Duration `mapstructure:"PRODUCER_LINGER"`
	DLQTopic       string        `mapstructure:"DLQ_TOPIC"`
}
