package config

import (
	"fmt"
	"net"
)

type DB struct {
	Host     string `mapstructure:"HOST"`
	Port     int    `mapstructure:"PORT"`
	User     string `mapstructure:"USER"`
	Password string `mapstructure:"PASSWORD"`
	DbName   string `mapstructure:"NAME"`
}

func (db *DB) GetPostgresURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		db.User, db.Password, net.JoinHostPort(db.Host, fmt.Sprint(db.Port)), db.DbName)
}
