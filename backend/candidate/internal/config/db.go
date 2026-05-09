package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type DB struct {
	Host     string `mapstructure:"HOST"`
	Port     int    `mapstructure:"PORT"`
	User     string `mapstructure:"USER"`
	Password string `mapstructure:"PASSWORD"`
	DbName   string `mapstructure:"NAME"`
}

func (db *DB) GetPostgresURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User, db.Password, db.Host, db.Port, db.DbName)
}

func bindEnvDB(v *viper.Viper) error {
	if err := v.BindEnv("DB.HOST", "DB_HOST"); err != nil {
		return err
	}
	if err := v.BindEnv("DB.PORT", "DB_PORT"); err != nil {
		return err
	}
	if err := v.BindEnv("DB.USER", "DB_USER"); err != nil {
		return err
	}
	if err := v.BindEnv("DB.PASSWORD", "DB_PASSWORD"); err != nil {
		return err
	}
	if err := v.BindEnv("DB.NAME", "DB_NAME"); err != nil {
		return err
	}
	return nil
}
