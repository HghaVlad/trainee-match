package config

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DB             DB             `mapstructure:"db"`
	HTTP           HTTP           `mapstructure:"http"`
	Kafka          Kafka          `mapstructure:"kafka"`
	SchemaRegistry SchemaRegistry `mapstructure:"schema_registry"`
	KafkaHandling  KafkaHandling  `mapstructure:"kafka_handling"`
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("db.max_pool_conns", 24)
	v.SetDefault("db.min_pool_conns", 2)

	_ = v.BindEnv("db.host")
	_ = v.BindEnv("db.port")
	_ = v.BindEnv("db.user")
	_ = v.BindEnv("db.password")
	_ = v.BindEnv("db.name")
	_ = v.BindEnv("http.addr")
	_ = v.BindEnv("http.jwkurl")

	_ = v.BindEnv("kafka.BROKERS")
	_ = v.BindEnv("kafka.CLIENT_ID")
	_ = v.BindEnv("kafka.CONSUMER_GROUP")
	_ = v.BindEnv("kafka.CONSUMER_TOPICS")
	_ = v.BindEnv("kafka.PRODUCER_ACKS")
	_ = v.BindEnv("kafka.PRODUCER_LINGER")
	_ = v.BindEnv("schema_registry.base_url")

	_ = v.BindEnv("kafka_handling.retry_delay")
	_ = v.BindEnv("kafka_handling.retry_count")

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	validate := validator.New()

	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
