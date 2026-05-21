package config

import (
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	Db             DB             `mapstructure:"DB"`
	Kafka          Kafka          `mapstructure:"KAFKA"`
	Outbox         Outbox         `mapstructure:"OUTBOX"`
	SchemaRegistry SchemaRegistry `mapstructure:"SCHEMA_REGISTRY"`
	Addr           string         `mapstructure:"ADDR"`
	JWKUrl         string         `mapstructure:"JWKURL"`
	GrpcAddr       string         `mapstructure:"GRPCADDR"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// bind envs grouped by sections
	if err := bindEnvDB(v); err != nil {
		return nil, err
	}
	if err := bindAddrAndAuthEnv(v); err != nil {
		return nil, err
	}
	if err := bindKafkaEnv(v); err != nil {
		return nil, err
	}
	if err := bindOutboxEnv(v); err != nil {
		return nil, err
	}
	if err := bindSchemaRegistryEnv(v); err != nil {
		return nil, err
	}

	readConfigFile(v)

	cfg, err := unmarshalConfig(v)
	if err != nil {
		return nil, err
	}

	parseBrokers(v, cfg)

	return cfg, nil
}

func bindAddrAndAuthEnv(v *viper.Viper) error {
	if err := v.BindEnv("Addr", "ADDR"); err != nil {
		return err
	}
	if err := v.BindEnv("JWKUrl", "JWKURL"); err != nil {
		return err
	}
	if err := v.BindEnv("GrpcAddr", "GRPCADDR"); err != nil {
		return err
	}
	return nil
}

// readConfigFile attempts to read a local config.env file (if present)
func readConfigFile(v *viper.Viper) {
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")

}

// unmarshalConfig unmarshals viper settings into Config with decode hooks
func unmarshalConfig(v *viper.Viper) (*Config, error) {
	var cfg Config
	if err := v.Unmarshal(&cfg, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
	))); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// parseBrokers parses comma separated brokers env into cfg.Kafka.Brokers
func parseBrokers(v *viper.Viper, cfg *Config) {
	if v.IsSet("KAFKA.BROKERS") {
		raw := v.GetString("KAFKA.BROKERS")
		if raw != "" {
			parts := strings.Split(raw, ",")
			for i := range parts {
				parts[i] = strings.TrimSpace(parts[i])
			}
			cfg.Kafka.Brokers = parts
		}
	}
}
