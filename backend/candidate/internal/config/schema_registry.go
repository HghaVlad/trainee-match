package config

import "github.com/spf13/viper"

type SchemaRegistry struct {
	BaseURL string `mapstructure:"BASE_URL"`
}

// bindSchemaRegistryEnv binds schema registry env variables
func bindSchemaRegistryEnv(v *viper.Viper) error {
	if err := v.BindEnv("SCHEMA_REGISTRY.BASE_URL", "SCHEMA_REGISTRY_BASE_URL"); err != nil {
		return err
	}
	return nil
}
