package config

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DB   DB   `mapstructure:"db"`
	Http HTTP `mapstructure:"http"`
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
