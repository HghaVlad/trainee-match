package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	KeyCloack KeyCloack `mapstructure:"KC"`
	Addr      string    `mapstructure:"ADDR"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("KC.URL", "http://0.0.0.0:800")
	v.SetDefault("KC.REALM", "trainee-match")
	v.SetDefault("KC.CLIENT_ID", "auth_backend")
	v.SetDefault("KC.CLIENT_SECRET", "UAaVKlGqGXZs2LXZPiV3uFYCblmrEhJ8")
	v.SetDefault("KC.ADMIN_USERNAME", "admin")
	v.SetDefault("KC.ADMIN_PASSWORD", "admin")

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
