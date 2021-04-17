package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"
)

// Connect to SQL database specified in configuration.
func Connect(cfg *Config) (*sql.DB, error) {
	dsn, err := url.Parse(cfg.DSN)

	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	db, err := sql.Open(cfg.Driver, dsn.String())

	if err != nil {
		return nil, fmt.Errorf("failed to open SQL connection: %w", err)
	}

	db.SetMaxOpenConns(cfg.Connections.Max)
	db.SetMaxIdleConns(cfg.Connections.Idle.Count)
	db.SetConnMaxLifetime(time.Duration(cfg.Connections.Idle.Time) * time.Second)

	return db, nil
}
