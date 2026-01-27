package infra_redis

import "github.com/HghaVlad/trainee-match/backend/company/internal/config"

type Config struct {
	Addr string
}

func NewConfig(conf *config.RedisConfig) *Config {
	return &Config{
		Addr: conf.Host + ":" + conf.Port,
	}
}
