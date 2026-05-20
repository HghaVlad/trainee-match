package config

import "time"

type Kafka struct {
	Brokers  []string `mapstructure:"brokers"`
	ClientID string   `mapstructure:"client_id"`

	ConsumerGroup  string   `mapstructure:"consumer_group"`
	ConsumerTopics []string `mapstructure:"consumer_topics"`

	ProducerAcks   string        `mapstructure:"producer_acks"`
	ProducerLinger time.Duration `mapstructure:"producer_linger"`
	DLQTopic       string        `mapstructure:"dlq_topic"`
}
