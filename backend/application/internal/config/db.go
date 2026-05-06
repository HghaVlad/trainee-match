package config

import (
	"fmt"
	"net"
)

type DB struct {
	Host         string `mapstructure:"host"           validate:"required"`
	Port         string `mapstructure:"port"           validate:"required"`
	User         string `mapstructure:"user"           validate:"required"`
	Password     string `mapstructure:"password"       validate:"required"`
	DbName       string `mapstructure:"name"           validate:"required"`
	MaxPoolConns int    `mapstructure:"max_pool_conns" validate:"gte=1,lte=100"`
	MinPoolConns int    `mapstructure:"min_pool_conns" validate:"gte=0,lte=100"`
}

func (db *DB) GetPostgresURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		db.User, db.Password, net.JoinHostPort(db.Host, db.Port), db.DbName)
}
