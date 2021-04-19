package provider

import (
	"context"

	"{{[ .Project ]}}/contracts/events"
)

// Events defines data store Events provider methods.
type Events interface {
	TransactProvider() (EventsTransact, error)
	Context(ctx context.Context) Events
	Create(model *events.Item) (*events.Item, error)
	Find(id string) (*events.Item, error)
	FindByName(name string, pageParams ...interface{}) ([]*events.Item, error)
	List(pageParams ...interface{}) ([]*events.Item, error)
	Update(model *events.Item) (*events.Item, error)
	Delete(id string) error
	DeleteByName(name string) error
}

// EventsTransact allow transactions in provider.
type EventsTransact interface {
	Transact
	Events
}
