package config

type BrokerConfig struct {
	Brokers       []string
	ConsumerGroup string
	// smth like OrderEventsTopic            string
}
