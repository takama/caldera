package postgres

import (
	"context"
	"database/sql"

	"{{[ .Project ]}}/pkg/db"
	"{{[ .Project ]}}/pkg/db/migrations"
	{{[- if .Contract ]}}
	"{{[ .Project ]}}/pkg/db/provider"
	{{[- end ]}}

	// Postgres driver
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	// Driver defines database driver name
	Driver = "postgres"
)

// Postgres controls postgres driver connection and providers
type Postgres struct {
	pool *sql.DB
	cfg  *db.Config
	log  *zap.Logger
	{{[- if .Contract ]}}
	// Contract providers
	events provider.Events
	{{[- end ]}}
}

// New creates new postgres DB connection
func New(cfg *db.Config, log *zap.Logger, mig migrations.Migrator) (*Postgres, error) {
	p := &Postgres{
		cfg: cfg,
		log: log,
	}
	var err error
	p.pool, err = db.Connect(cfg)
	if err != nil {
		return nil, err
	}
	name := cfg.Driver
	if err := p.pool.QueryRow("SELECT version()").Scan(&name); err != nil {
		return nil, err
	}

	log.Info("DB", zap.String("version", name))

	{{[- if .Contract ]}}

	p.events = newEventsProvider(p.pool)
	{{[- end ]}}

	// setup migration connection
	mig.Setup(p.pool)

	return p, mig.Migrate()
}

// Check readiness for database
func (p Postgres) Check() error {
	return p.pool.Ping()
}

// Shutdown process graceful shutdown for the server
func (p Postgres) Shutdown(ctx context.Context) error {
	return p.pool.Close()
}

{{[- if .Contract ]}}

// EventsProvider returns data store provider for Events
func (p Postgres) EventsProvider() provider.Events {
	return p.events
}
{{[- end ]}}
