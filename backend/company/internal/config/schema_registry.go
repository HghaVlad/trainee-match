package config

import "github.com/go-playground/validator/v10"

type SchemaRegistry struct {
	BaseURL string `env:"SCHEMA_REGISTRY_BASE_URL" validate:"required,url"`
}

func LoadSchemaRegistryConfig(validate *validator.Validate) (*SchemaRegistry, error) {
	return loadConfigSection[SchemaRegistry](validate, "schema_registry")
}
