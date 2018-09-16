package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
)

// Connect to SQL database specified in configuration
func Connect(cfg *Config) (*sql.DB, error) {
	var properties string
	if len(cfg.Properties) > 0 {
		properties = "?" + strings.Join(cfg.Properties, "&")
	}
	dsn, err := url.Parse(fmt.Sprintf("%s://%s:%s@%s:%d/%s%s",
		cfg.Driver, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name, properties))
	if err != nil {
		return nil, err
	}
	db, err := sql.Open(cfg.Driver, dsn.String())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.Connections.Max)
	db.SetMaxIdleConns(cfg.Connections.Idle)

	return db, nil
}
