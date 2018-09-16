package db

import (
	"context"

	"github.com/{{[ .Github ]}}/{{[ .Name ]}}/pkg/db/provider"
)

// Store design database interface with providers
type Store interface {
	Check() error
	Shutdown(ctx context.Context) error
	{{[- if .Contract ]}}
	EventsProvider() provider.Events
	{{[- end ]}}
}
