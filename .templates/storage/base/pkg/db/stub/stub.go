package stub

import (
	"context"
	"os"

	"{{[ .Project ]}}/pkg/db"
	{{[- if .Contract ]}}
	"{{[ .Project ]}}/pkg/db/provider"
	{{[- end ]}}

	"go.uber.org/zap"
)

const (
	// Driver defines database driver name
	Driver = "stub"
)

// Stub contains stub configuration
type Stub struct {
	cfg *db.Config
	log *zap.Logger
	{{[- if .Contract ]}}
	// Contract providers
	events *eventsProvider
	{{[- end ]}}
}

// New creates new DB connection
func New(cfg *db.Config, log *zap.Logger) (*Stub, error) {
	log.Info("DB", zap.String("driver", Driver))

	conn := &Stub{
		cfg: cfg,
		log: log,
	}

	{{[- if .Contract ]}}

	conn.events = &eventsProvider{cfg: cfg}

	if err := conn.events.load(); err != nil {
		return nil, err
	}
	{{[- end ]}}

	return conn, nil
}

// Check DB
func (s Stub) Check() error {
	return nil
}

// Shutdown process graceful shutdown for the storage
func (s Stub) Shutdown(ctx context.Context) error {
	s.log.Debug("DB closed", zap.String("driver", Driver))
	return nil
}

{{[- if .Contract ]}}

// EventsProvider returns data store provider for Events
func (s Stub) EventsProvider() provider.Events {
	return s.events
}
{{[- end ]}}

func readFile(path string) (*os.File, error) {
	_, err := os.Stat(path)
	// if file does not exist, return "empty data" without error
	if os.IsNotExist(err) {
		return nil, nil
	}
	return os.Open(path) // nolint: gosec
}
