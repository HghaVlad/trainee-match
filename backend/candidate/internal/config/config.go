package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Db     DB     `mapstructure:"DB"`
	Addr   string `mapstructure:"ADDR"`
	JWKUrl string `mapstructure:"JWKURL"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("DB.HOST", "localhost")
	v.SetDefault("DB.PORT", 5432)
	v.SetDefault("DB.USER", "postgres")
	v.SetDefault("DB.PASSWORD", "postgres")
	v.SetDefault("DB.NAME", "candidate")

	v.BindEnv("Addr", "ADDR")
	v.BindEnv("JWKUrl", "JWKURL")

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
