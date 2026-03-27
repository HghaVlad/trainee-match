package config

type Postgres struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string

	MaxPoolConns int
	MinPoolConns int
}
