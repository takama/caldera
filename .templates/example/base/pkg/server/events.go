package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"{{[ .Project ]}}/contracts/events"
	"{{[ .Project ]}}/contracts/request"
	"{{[ .Project ]}}/pkg/db/provider"
)

type eventsServer struct {
	provider.Events
}

// GetEvent returns Event requested by ID
func (es eventsServer) GetEvent(
	ctx context.Context,
	req *request.ByID,
) (*events.Event, error) {
	return es.Context(ctx).Find(req.Id)
}

// FindEventsByName returns Event objects requested by event name
func (es eventsServer) FindEventsByName(
	req *request.ByName,
	stream events.Events_FindEventsByNameServer,
) error {
	data, err := es.FindByName(req.Name)
	if err != nil {
		return err
	}
	for ind := range data {
		if err := stream.Send(&data[ind]); err != nil {
			return err
		}
	}
	return nil
}

// ListEvents returns all Event objects
func (es eventsServer) ListEvents(
	empty *empty.Empty,
	stream events.Events_ListEventsServer,
) error {
	data, err := es.List()
	if err != nil {
		return err
	}
	for ind := range data {
		if err := stream.Send(&data[ind]); err != nil {
			return err
		}
	}
	return nil
}

// CreateEvent creates a new Event object
func (es eventsServer) CreateEvent(
	ctx context.Context,
	event *events.Event,
) (*events.Event, error) {
	return es.Context(ctx).Create(event)
}

// UpdateEvent updates an existing Event object
func (es eventsServer) UpdateEvent(
	ctx context.Context,
	event *events.Event,
) (*events.Event, error) {
	return es.Context(ctx).Update(event)
}

// DeleteEvent removes Event object requested by ID
func (es eventsServer) DeleteEvent(
	ctx context.Context,
	req *request.ByID,
) (*empty.Empty, error) {
	return new(empty.Empty), es.Context(ctx).Delete(req.Id)
}

// DeleteEventsByName removes Event objects requested by event name
func (es eventsServer) DeleteEventsByName(
	ctx context.Context,
	req *request.ByName,
) (*empty.Empty, error) {
	return new(empty.Empty), es.Context(ctx).DeleteByName(req.Name)
}
