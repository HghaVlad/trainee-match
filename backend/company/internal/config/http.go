package config

import "github.com/go-playground/validator/v10"

type HTTP struct {
	Addr   string `env:"HTTP_ADDR" envDefault:":8080" validate:"required"`
	JWKUrl string `env:"JWK_URL"                      validate:"required"`
}

func LoadHTTPConfig(validate *validator.Validate) (*HTTP, error) {
	return loadConfigSection[HTTP](validate, "http")
}
