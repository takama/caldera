package mysql

import (
	"context"
	"database/sql"

	"{{[ .Project ]}}/pkg/db"
	"{{[ .Project ]}}/pkg/db/migrations"
	{{[- if .Contract ]}}
	"{{[ .Project ]}}/pkg/db/provider"
	{{[- end ]}}

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

const (
	// Driver defines database driver name
	Driver = "mysql"
)

// MySQL controls mysql driver connection and providers
type MySQL struct {
	pool *sql.DB
	cfg  *db.Config
	log  *zap.Logger
	{{[- if .Contract ]}}
	// Contract providers
	events provider.Events
	{{[- end ]}}
}

// New creates new postgres DB connection
func New(cfg *db.Config, log *zap.Logger, mig migrations.Migrator) (*MySQL, error) {
	m := &MySQL{
		cfg: cfg,
		log: log,
	}
	var err error
	m.pool, err = db.Connect(cfg)
	if err != nil {
		return nil, err
	}
	name := cfg.Driver
	if err := m.pool.QueryRow("SELECT version()").Scan(&name); err != nil {
		return nil, err
	}

	log.Info("DB", zap.String("version", name))

	{{[- if .Contract ]}}

	m.events = newEventsProvider(m.pool)
	{{[- end ]}}

	// setup migration connection
	mig.Setup(m.pool)

	return m, mig.Migrate()
}

// Check readiness for database
func (m MySQL) Check() error {
	return m.pool.Ping()
}

// Shutdown process graceful shutdown for the server
func (m MySQL) Shutdown(ctx context.Context) error {
	return m.pool.Close()
}

{{[- if .Contract ]}}

// EventsProvider returns data store provider for Events
func (m MySQL) EventsProvider() provider.Events {
	return m.events
}
{{[- end ]}}
