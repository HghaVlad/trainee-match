package config

import "time"

type Kafka struct {
	Brokers  []string
	ClientId string

	ProducerAcks   string
	ProducerLinger time.Duration
}
