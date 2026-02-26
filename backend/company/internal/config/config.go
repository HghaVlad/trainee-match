package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTP      HTTPConfig
	CompanyDB DBConfig
	Broker    BrokerConfig
	Redis     RedisConfig
}

func Load() (*Config, error) {
	cfg := &Config{
		HTTP: HTTPConfig{
			Addr:   getEnv("HTTP_ADDR", ":8088"),
			JWKUrl: getEnv("JWK_URL", ""),
		},
		CompanyDB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),

			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 30*time.Second),
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

func getEnvDuration(key string, def time.Duration) time.Duration {
	if str := os.Getenv(key); str != "" {
		num, err := strconv.Atoi(str)
		if err != nil {
			return def
		}
		return time.Duration(num) * time.Second
	}

	return def
}
