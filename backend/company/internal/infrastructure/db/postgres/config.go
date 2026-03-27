package postgres

import (
	"time"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
)

type Config struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Host:     cfg.CompanyDB.Host,
		Port:     cfg.CompanyDB.Port,
		Name:     cfg.CompanyDB.Name,
		User:     cfg.CompanyDB.User,
		Password: cfg.CompanyDB.Password,
		SSLMode:  cfg.CompanyDB.SSLMode,

		MaxOpenConns:    cfg.CompanyDB.MaxOpenConns,
		MaxIdleConns:    cfg.CompanyDB.MaxIdleConns,
		ConnMaxLifetime: cfg.CompanyDB.ConnMaxLifetime,
	}
}
