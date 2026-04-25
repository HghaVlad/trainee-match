package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HTTP     HTTPConfig
	Postgres Postgres
	Broker   BrokerConfig
	Redis    RedisConfig
}

func Load() (*Config, error) {
	cfg := &Config{
		HTTP: HTTPConfig{
			Addr:   getEnv("HTTP_ADDR", ":8088"),
			JWKUrl: getEnv("JWK_URL", ""),
		},
		Postgres: Postgres{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Name:     os.Getenv("POSTGRES_DB"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),

			MaxPoolConns: getEnvInt("POSTGRES_MAX_POOL_CONNS", 10),
			MinPoolConns: getEnvInt("POSTGRES_MIN_POOL_CONNS", 2),
		},
		Broker: BrokerConfig{
			Brokers:       strings.Split(os.Getenv("BROKER_HOST"), ","),
			ConsumerGroup: os.Getenv("BROKER_CONSUMER_GROUP"),
		},
		Redis: RedisConfig{
			Host: os.Getenv("REDIS_HOST"),
			Port: os.Getenv("REDIS_PORT"),
		},
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if str := os.Getenv(key); str != "" {
		num, err := strconv.Atoi(str)
		if err != nil {
			return def
		}
		return num
	}
	return def
}
