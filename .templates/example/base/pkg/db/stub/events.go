package stub

import (
	"bufio"
	"context"
	"encoding/json"
	"path/filepath"
	"sync"

	"{{[ .Project ]}}/contracts/events"
	"{{[ .Project ]}}/pkg/db"
	"{{[ .Project ]}}/pkg/db/provider"
)

type eventsProvider struct {
	mutex sync.RWMutex
	cfg   *db.Config
	Data  []events.Event
}

// Transaction returns provider with transaction
func (ep *eventsProvider) TransactProvider() (provider.EventsTransact, error) {
	return ep, nil
}

// Commit changes in depth of transaction
func (ep *eventsProvider) Commit() error {
	return nil
}

// Rollback changes in depth of transaction
func (ep *eventsProvider) Rollback() error {
	return nil
}

// Context returns provider with context
func (ep *eventsProvider) Context(ctx context.Context) provider.Events {
	return ep
}

// Create new Event object
func (ep *eventsProvider) Create(model *events.Event) (*events.Event, error) {
	ep.mutex.Lock()
	defer ep.mutex.Unlock()

	if model.Id != "" {
		ind, item := ep.findByID(model.Id)
		if ind == -1 {
			return nil, provider.ErrNotExistingEvent
		}
		if item != nil {
			return nil, provider.ErrAlreadyExistingID
		}
	}
	ep.Data = append(ep.Data, *model)
	return model, nil
}

// Find returns Event requested by ID
func (ep *eventsProvider) Find(id string) (*events.Event, error) {
	ind, item := ep.findByID(id)
	if ind == -1 {
		return nil, provider.ErrNotExistingEvent
	}

	return item, nil
}

// FindByName returns Events requested by Event name
func (ep *eventsProvider) FindByName(name string) ([]events.Event, error) {
	_, items := ep.findByName(name)
	return items, nil
}

// List returns all Event objects
func (ep *eventsProvider) List() ([]events.Event, error) {
	items := make([]events.Event, len(ep.Data))
	ep.mutex.RLock()
	defer ep.mutex.RUnlock()
	copy(items, ep.Data)
	return items, nil
}

// Update Event object
func (ep *eventsProvider) Update(model *events.Event) (*events.Event, error) {
	ind, _ := ep.findByID(model.Id)
	if ind == -1 {
		return nil, provider.ErrNotExistingEvent
	}

	ep.mutex.Lock()
	defer ep.mutex.Unlock()
	ep.Data = append(append(ep.Data[:ind], *model), ep.Data[ind+1:]...)
	return model, nil
}

// Delete removes Event object by ID
func (ep *eventsProvider) Delete(id string) error {
	ind, _ := ep.findByID(id)
	if ind == -1 {
		return provider.ErrNotExistingEvent
	}

	ep.mutex.Lock()
	defer ep.mutex.Unlock()
	ep.Data = append(ep.Data[:ind], ep.Data[ind+1:]...)
	return nil
}

// DeleteByName removes Event objects by Event name
func (ep *eventsProvider) DeleteByName(name string) error {
	indices, _ := ep.findByName(name)
	if len(indices) == 0 {
		return nil
	}

	for len(indices) != 0 {
		ep.mutex.Lock()
		ep.Data = append(ep.Data[:indices[0]], ep.Data[indices[0]+1:]...)
		ep.mutex.Unlock()
		indices, _ = ep.findByName(name)
	}
	return nil
}

func (ep *eventsProvider) findByID(id string) (int, *events.Event) {
	ep.mutex.RLock()
	defer ep.mutex.RUnlock()

	for ind := range ep.Data {
		if ep.Data[ind].Id == id {
			return ind, &ep.Data[ind]
		}
	}

	return -1, nil
}

func (ep *eventsProvider) findByName(name string) ([]int, []events.Event) {
	indices := make([]int, 0)
	items := make([]events.Event, 0)
	ep.mutex.RLock()
	defer ep.mutex.RUnlock()

	for k, v := range ep.Data {
		if v.Name == name {
			indices = append(indices, k)
			items = append(items, v)
		}
	}

	return indices, items
}

func (ep *eventsProvider) load() error {
	ep.Data = make([]events.Event, 0)
	path := filepath.Join(ep.cfg.Fixtures.Dir, "events/data.json")
	f, err := readFile(path)
	if err != nil || f == nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(bufio.NewReader(f)).Decode(&ep)
}
