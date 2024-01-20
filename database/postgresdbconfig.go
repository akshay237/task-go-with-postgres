package database

import (
	"fmt"
)

type ConfigPostgres struct {
	Host             string
	User             string
	Password         string
	Database         string
	SSLMode          string
	ConnMaxOpen      int
	ConnMaxIdleTime  int64
	ConnMaxIdleConns int
}

func (cfg *ConfigPostgres) FormatDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=verify-full",
		cfg.User, cfg.Password, cfg.Host, cfg.Database)
}
