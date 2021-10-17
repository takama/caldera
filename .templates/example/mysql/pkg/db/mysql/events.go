package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"{{[ .Project ]}}/contracts/events"
	"{{[ .Project ]}}/pkg/db/provider"

	"github.com/satori/go.uuid"
)

type eventsProvider struct {
	*provider.SQL
}

func newEventsProvider(db *sql.DB) *eventsProvider {
	return &eventsProvider{SQL: provider.New(db)}
}

// Transaction returns provider with transaction.
func (ep *eventsProvider) TransactProvider() (provider.EventsTransact, error) {
	p, err := ep.SQL.TransactProvider()
	if err != nil {
		return ep, fmt.Errorf("failed to create transact provider: %w", err)
	}

	return &eventsProvider{SQL: p}, nil
}

// Context returns provider with context.
func (ep *eventsProvider) Context(ctx context.Context) provider.Events {
	return &eventsProvider{SQL: ep.SQL.Context(ctx)}
}

// Create new Event object.
func (ep *eventsProvider) Create(model *events.Item) (*events.Item, error) {
	if model.Name == "" {
		return nil, provider.ErrNotDefinedName
	}

	model.Id = uuid.NewV4().String()
	stmt, err := ep.Prepare(queryInsertEvent)

	if err != nil {
		return nil, fmt.Errorf("failed to prepare create event request: %w", err)
	}

	defer stmt.Close()

	if _, err = stmt.Exec(model.Id, model.Name); err != nil {
		return model, fmt.Errorf("failed to create event: %w", err)
	}

	return model, nil
}

// Find returns Event requested by ID.
func (ep *eventsProvider) Find(id string) (*events.Item, error) {
	event := new(events.Item)
	row := ep.QueryRow(queryEventByID, id)

	if err := row.Scan(&event.Id, &event.Name); err != nil {
		return event, fmt.Errorf("failed to find event(s): %w", err)
	}

	return event, nil
}

// FindByName returns Events requested by Event name.
func (ep *eventsProvider) FindByName(name string, pageParams ...interface{}) ([]*events.Item, error) {
	params := make([]interface{}, 0)
	params = append(params, name)
	params = append(params, pageParams...)

	return ep.find(queryEventsByName, params...)
}

// List returns all Event objects.
func (ep *eventsProvider) List(pageParams ...interface{}) ([]*events.Item, error) {
	return ep.find(queryEvents, pageParams...)
}

// Update Event object.
func (ep *eventsProvider) Update(model *events.Item) (*events.Item, error) {
	if model.Id == "" {
		return nil, provider.ErrNotDefinedID
	}

	if model.Name == "" {
		return nil, provider.ErrNotDefinedName
	}

	stmt, err := ep.Prepare(queryUpdateEvent)
	
	if err != nil {
		return nil, fmt.Errorf("failed to prepare update event request: %w", err)
	}

	defer stmt.Close()

	if _, err = stmt.Exec(model.Id, model.Name); err != nil {
		return model, fmt.Errorf("failed to update event: %w", err)
	}

	return model, nil
}

// Delete removes Event object by ID.
func (ep *eventsProvider) Delete(id string) error {
	if id == "" {
		return provider.ErrNotDefinedID
	}

	stmt, err := ep.Prepare(queryDeleteEventByID)
	
	if err != nil {
		return fmt.Errorf("failed to prepare delete event request: %w", err)
	}
	
	defer stmt.Close()	
	_, err = stmt.Exec(id)

	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

// DeleteByName removes Event objects by Event name.
func (ep *eventsProvider) DeleteByName(name string) error {
	if name == "" {
		return provider.ErrNotDefinedName
	}

	stmt, err := ep.Prepare(queryDeleteEventsByName)

	if err != nil {
		return fmt.Errorf("failed to prepare delete event by name request: %w", err)
	}

	defer stmt.Close()
	_, err = stmt.Exec(name)

	if err != nil {
		return fmt.Errorf("failed to delete event by name: %w", err)
	}

	return nil
}

func (ep *eventsProvider) find(query string, args ...interface{}) ([]*events.Item, error) {
	items := make([]*events.Item, 0)
	rows, err := ep.Query(query, args...)

	if err != nil {
		return nil, fmt.Errorf("failed to prepare find: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		item := new(events.Item)
		if err := rows.Scan(&item.Id, &item.Name); err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return items, fmt.Errorf("failed to find event(s): %w", rows.Err())
	}

	return items, nil
}

const (
	queryEventByID    = `SELECT id, name FROM events WHERE id = ?`
	queryEventsByName = `SELECT id, name FROM events WHERE name = ? LIMIT ? OFFSET ?`
	queryEvents       = `SELECT id, name FROM events LIMIT ? OFFSET ?`
	queryInsertEvent  = `INSERT INTO events (id, name) VALUES (?, ?)`
	queryUpdateEvent  = `INSERT INTO events (id, name) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE name = ?`
	queryDeleteEventByID    = `DELETE FROM events WHERE id = ?`
	queryDeleteEventsByName = `DELETE FROM events WHERE name = ?`
)
