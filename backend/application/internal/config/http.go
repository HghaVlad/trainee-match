package config

type Http struct {
	Addr   string `mapstructure:"addr"`
	JWKUrl string `mapstructure:"jwkurl"`
}
