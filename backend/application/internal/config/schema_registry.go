package config

import "time"

type SchemaRegistry struct {
	BaseURL string
	TimeOut time.Duration
}
