package provider

import (
	"context"

	"{{[ .Project ]}}/contracts/events"
)

// Events defines data store Events provider methods
type Events interface {
	TransactProvider() (EventsTransact, error)
	Context(ctx context.Context) Events
	Create(model *events.Event) (*events.Event, error)
	Find(id string) (*events.Event, error)
	FindByName(name string) ([]events.Event, error)
	List() ([]events.Event, error)
	Update(model *events.Event) (*events.Event, error)
	Delete(id string) error
	DeleteByName(name string) error
}

// EventsTransact allow transactions in provider
type EventsTransact interface {
	Transact
	Events
}
