package database

import (
	"fmt"
)

type ConfigPostgres struct {
	Host             string
	Port             int
	User             string
	Password         string
	Database         string
	SSLMode          string
	ConnMaxOpen      int
	ConnMaxIdleTime  int64
	ConnMaxIdleConns int
}

func (cfg *ConfigPostgres) FormatDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Database, cfg.Password)
}
