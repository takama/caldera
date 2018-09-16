package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/{{[ .Github ]}}/{{[ .Name ]}}/contracts/events"
	"github.com/{{[ .Github ]}}/{{[ .Name ]}}/contracts/request"
	"github.com/{{[ .Github ]}}/{{[ .Name ]}}/pkg/db/provider"
)

type eventsServer struct {
	provider.Events
}

// GetEvent returns Event requested by ID
func (es eventsServer) GetEvent(
	ctx context.Context,
	req *request.ByID,
) (*events.Event, error) {
	return es.Context(ctx).Find(req.ID)
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
	for _, r := range data {
		if err := stream.Send(&r); err != nil {
			return err
		}
	}
	return nil
}

// ListEvents returns all Branch objects
func (es eventsServer) ListEvents(
	empty *empty.Empty,
	stream events.Events_ListEventsServer,
) error {
	data, err := es.List()
	if err != nil {
		return err
	}
	for _, r := range data {
		if err := stream.Send(&r); err != nil {
			return err
		}
	}
	return nil
}

// NewEvent creates a new Event object
func (es eventsServer) NewEvent(
	ctx context.Context,
	event *events.Event,
) (*events.Event, error) {
	return es.Context(ctx).New(event)
}

// SaveEvent updates an existing Event object
func (es eventsServer) SaveEvent(
	ctx context.Context,
	event *events.Event,
) (*events.Event, error) {
	return es.Context(ctx).Save(event)
}

// DeleteEvent removes Event object requested by ID
func (es eventsServer) DeleteEvent(
	ctx context.Context,
	req *request.ByID,
) (*empty.Empty, error) {
	return new(empty.Empty), es.Context(ctx).Delete(req.ID)
}

// DeleteEventsByName removes Event objects requested by event name
func (es eventsServer) DeleteEventsByName(
	ctx context.Context,
	req *request.ByName,
) (*empty.Empty, error) {
	return new(empty.Empty), es.Context(ctx).DeleteByName(req.Name)
}
