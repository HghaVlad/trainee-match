package config

type HTTP struct {
	Addr   string `mapstructure:"addr"`
	JWKUrl string `mapstructure:"jwkurl"`
}
