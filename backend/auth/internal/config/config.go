package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	KeyCloack KeyCloack `mapstructure:"KC"`
	Addr      string    `mapstructure:"ADDR"`
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

	v.SetConfigName("config")
	v.SetConfigType("env")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err == nil {
		fmt.Printf("Found file %s. Using config from file\n", v.ConfigFileUsed())
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
