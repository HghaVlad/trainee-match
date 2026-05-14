package config

import "time"

type SchemaRegistry struct {
	BaseURL string        `mapstructure:"base_url"`
	TimeOut time.Duration `mapstructure:"timeout"`
}
