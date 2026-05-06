package config

type DB struct {
	Host         string `mapstructure:"host"           validate:"required"`
	Port         string `mapstructure:"port"           validate:"required"`
	User         string `mapstructure:"user"           validate:"required"`
	Password     string `mapstructure:"password"       validate:"required"`
	DBName       string `mapstructure:"name"           validate:"required"`
	MaxPoolConns int    `mapstructure:"max_pool_conns" validate:"gte=1,lte=100"`
	MinPoolConns int    `mapstructure:"min_pool_conns" validate:"gte=0,lte=100"`
}

