package config

import (
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Db     DB     `mapstructure:"DB"`
	Addr   string `mapstructure:"ADDR"`
	JWKUrl string `mapstructure:"JWKURL"`
}

func Load() *Config {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.BindEnv("DB.HOST")
	v.BindEnv("DB.PORT")
	v.BindEnv("DB.USER")
	v.BindEnv("DB.PASSWORD")
	v.BindEnv("DB.NAME")

	v.BindEnv("ADDR")
	v.BindEnv("JWKURL")

	v.AutomaticEnv()

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		slog.Error("Error parsing config file, %s", err)
	}

	return &config
}
