package config

import "github.com/go-playground/validator/v10"

type Redis struct {
	Host string `env:"REDIS_HOST" validate:"required"`
	Port string `env:"REDIS_PORT" validate:"required,numeric"`
}

func LoadRedisConfig(validate *validator.Validate) (*Redis, error) {
	return loadConfigSection[Redis](validate, "redis")
}
