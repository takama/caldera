package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"{{[ .Project ]}}/pkg/db"
	{{[- if .Prometheus.Enabled ]}}
	"{{[ .Project ]}}/pkg/metrics"
	{{[- end ]}}
	"{{[ .Project ]}}/pkg/db/migrations"
	{{[- if .Example ]}}
	"{{[ .Project ]}}/pkg/db/provider"
	{{[- end ]}}

	// Postgres driver.
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	// Driver defines database driver name.
	Driver = "postgres"
)

// Postgres controls postgres driver connection and providers.
type Postgres struct {
	pool *sql.DB
	cfg  *db.Config
	log  *zap.Logger
	{{[- if .Example ]}}
	// Contract providers
	events provider.Events
	{{[- end ]}}
}

// DSN creates dsn type connection.
func DSN(cfg *db.Config) *db.Config {
	if cfg.DSN == "" {
		var properties string

		if len(cfg.Properties) > 0 {
			properties = "?" + strings.Join(cfg.Properties, "&")
		}

		dsn := fmt.Sprintf("%s://%s:%s@%s:%d/%s%s",
			cfg.Driver, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name, properties)
		cfg.DSN = dsn
	}

	return cfg
}

// New creates new postgres DB connection.
func New(cfg *db.Config, log *zap.Logger, mig migrations.Migrator) (*Postgres, error) {
	p := &Postgres{
		cfg: cfg,
		log: log,
	}

	var err error

	p.pool, err = db.Connect(cfg)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	name := cfg.Driver

	if err := p.pool.QueryRow("SELECT version()").Scan(&name); err != nil {
		return nil, fmt.Errorf("failed to check the database engine version: %w", err)
	}

	log.Info("DB", zap.String("version", name))

	{{[- if .Example ]}}

	p.events = newEventsProvider(p.pool)
	{{[- end ]}}

	// setup migration connection
	if err := mig.Setup(p.pool); err != nil {
		return p, fmt.Errorf("failed to setup migration %w", err)
	}

	if err := mig.Migrate(); err != nil {
		return p, fmt.Errorf("failed to init connection with migration %w", err)
	}

	return p, nil
}

// Check readiness for database.
func (p Postgres) Check() error {
	if err := p.pool.Ping(); err != nil {
		return fmt.Errorf("failed to check postgress connection %w", err)
	}

	return nil
}

// Shutdown process graceful shutdown for the server.
func (p Postgres) Shutdown(ctx context.Context) error {
	if err := p.pool.Close(); err != nil {
		return fmt.Errorf("failed to close postgress connection %w", err)
	}

	return nil
}

{{[- if .Example ]}}

// EventsProvider returns data store provider for Events.
func (p Postgres) EventsProvider() provider.Events {
	return p.events
}
{{[- end ]}}

{{[- if .Prometheus.Enabled ]}}

// MetricFunc returns a func to monitor connectivity to Postgres.
func (p Postgres) MetricFunc() metrics.MetricFunc {
	return metrics.DBMetricFunc(p.cfg.Host, p.cfg.Name, p.pool.Stats())
}
{{[- end ]}}
