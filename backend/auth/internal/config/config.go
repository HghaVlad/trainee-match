package config

import (
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	KeyCloak       KeyCloak       `mapstructure:"KC"`
	Addr           string         `mapstructure:"ADDR"`
	Kafka          Kafka          `mapstructure:"KAFKA"`
	SchemaRegistry SchemaRegistry `mapstructure:"SCHEMA_REGISTRY"`
	Outbox         Outbox         `mapstructure:"OUTBOX"`
	Postgres       Postgres       `mapstructure:"POSTGRES"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("KC.URL", "url")
	v.SetDefault("KC.REALM", "realm")
	v.SetDefault("KC.CLIENT_ID", "id")
	v.SetDefault("KC.CLIENT_SECRET", "secret")
	v.SetDefault("KC.ADMIN_USERNAME", "username")
	v.SetDefault("KC.ADMIN_PASSWORD", "password")

	v.SetDefault("KC.ACCESS_TOKEN_EXPIRES", 5*60)
	v.SetDefault("KC.REFRESH_TOKEN_EXPIRES", 30*60)

	v.SetDefault("ADDR", "0.0.0.0:8000")

	v.SetDefault("KAFKA.BROKERS", []string{"localhost:9092"})
	v.SetDefault("KAFKA.CLIENT_ID", "auth")
	v.SetDefault("KAFKA.PRODUCER_ACKS", "all")
	v.SetDefault("KAFKA.PRODUCER_LINGER", "10ms")
	v.SetDefault("KAFKA.USER_TOPIC", "user-created")

	v.SetDefault("SCHEMA_REGISTRY.URL", "http://localhost:8081")

	v.SetDefault("OUTBOX.USER_TOPIC", "user-created")
	v.SetDefault("OUTBOX.BASE_RETRY_DELAY", "5s")
	v.SetDefault("OUTBOX.MAX_RETRIES", 5)
	v.SetDefault("OUTBOX.BATCH_SIZE", 100)
	v.SetDefault("OUTBOX.RELAY_MIN_SLEEP", "50ms")
	v.SetDefault("OUTBOX.RELAY_MAX_SLEEP", "5s")
	v.SetDefault("OUTBOX.RELAY_WORKER_COUNT", 2)

	v.SetDefault("POSTGRES.HOST", "auth-postgres")
	v.SetDefault("POSTGRES.PORT", 5432)
	v.SetDefault("POSTGRES.USER", "postgres")
	v.SetDefault("POSTGRES.PASSWORD", "postgres")
	v.SetDefault("POSTGRES.DB_NAME", "auth_db")
	v.SetDefault("POSTGRES.SSL_MODE", "disable")

	v.SetConfigName("config")
	v.SetConfigType("env")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err == nil {
		slog.Debug("found file %s. Using config from file\n", v.ConfigFileUsed())
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
