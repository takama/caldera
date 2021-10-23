package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"

	"{{[ .Project ]}}/contracts/events"
	"{{[ .Project ]}}/pkg/db/provider"

	"go.uber.org/zap"
)

const (
	errorMessage = "Error"
	okMessage    = "Ok"
	maxPageSize  = 100
	numbersBase  = 10
	bitSize      = 32
)

type eventsServer struct {
	provider.Events
	log *zap.Logger
}

// GetEvent returns Event requested by ID.
func (es eventsServer) GetEvent(
	ctx context.Context,
	req *events.RequestByID,
) (*events.Response, error) {
	response := new(events.Response)
	value, err := es.Context(ctx).Find(req.Id)

	if err != nil {
		response.Message = errorMessage

		return response, fmt.Errorf("failed to find events by name: %w", err)
	}

	response.Message = okMessage
	response.Result = append(response.Result, value)
	response.Total = uint64(len(response.Result))

	return response, nil
}

// FindEventsByName returns Event objects requested by event name.
func (es eventsServer) FindEventsByName(
	ctx context.Context,
	req *events.RequestListByName,
) (*events.Response, error) {
	response := new(events.Response)

	limit, offset := es.offsetPaging(req.PageSize, req.PageToken)

	value, err := es.Context(ctx).FindByName(req.Name, limit, offset)

	if err != nil {
		response.Message = errorMessage

		return response, fmt.Errorf("failed to find events by name: %w", err)
	}

	response.Message = okMessage
	response.Result = append(response.Result, value...)
	response.Total = uint64(len(response.Result))

	if response.Total >= uint64(limit) {
		response.NextPageToken = es.offsetPageToken(limit, offset)
	} else {
		response.NextPageToken = req.PageToken
	}

	return response, nil
}

// ListEvents returns all Event objects.
func (es eventsServer) ListEvents(
	ctx context.Context,
	req *events.RequestList,
) (*events.Response, error) {
	response := new(events.Response)

	limit, offset := es.offsetPaging(req.PageSize, req.PageToken)

	value, err := es.Context(ctx).List(limit, offset)

	if err != nil {
		response.Message = errorMessage

		return response, fmt.Errorf("failed to List events: %w", err)
	}

	response.Message = okMessage
	response.Result = append(response.Result, value...)
	response.Total = uint64(len(response.Result))

	if response.Total >= uint64(limit) {
		response.NextPageToken = es.offsetPageToken(limit, offset)
	} else {
		response.NextPageToken = req.PageToken
	}

	return response, nil
}

// CreateEvent creates a new Event object.
func (es eventsServer) CreateEvent(
	ctx context.Context,
	event *events.Item,
) (*events.Response, error) {
	response := new(events.Response)

	value, err := es.Context(ctx).Create(event)

	if err != nil {
		response.Message = errorMessage

		return response, fmt.Errorf("failed to create events: %w", err)
	}

	response.Message = okMessage
	response.Result = append(response.Result, value)
	response.Total = uint64(len(response.Result))

	return response, nil
}

// UpdateEvent updates an existing Event object.
func (es eventsServer) UpdateEvent(
	ctx context.Context,
	event *events.Item,
) (*events.Response, error) {
	response := new(events.Response)

	value, err := es.Context(ctx).Update(event)

	if err != nil {
		response.Message = errorMessage

		return response, fmt.Errorf("failed to update events: %w", err)
	}

	response.Message = okMessage
	response.Result = append(response.Result, value)
	response.Total = uint64(len(response.Result))

	return response, nil
}

// DeleteEvent removes Event object requested by ID.
func (es eventsServer) DeleteEvent(
	ctx context.Context,
	req *events.RequestByID,
) (*events.Response, error) {
	response := new(events.Response)

	err := es.Context(ctx).Delete(req.Id)

	if err != nil {
		response.Message = errorMessage

		return response, fmt.Errorf("failed to delete events: %w", err)
	}

	response.Message = okMessage
	response.Total = uint64(len(response.Result))

	return response, nil
}

// DeleteEventsByName removes Event objects requested by event name.
func (es eventsServer) DeleteEventsByName(
	ctx context.Context,
	req *events.RequestByName,
) (*events.Response, error) {
	response := new(events.Response)

	err := es.Context(ctx).DeleteByName(req.Name)

	if err != nil {
		response.Message = errorMessage

		return response, fmt.Errorf("failed to delete events by name: %w", err)
	}

	response.Message = okMessage
	response.Total = uint64(len(response.Result))

	return response, nil
}

func (es eventsServer) offsetPaging(pageSize int32, pageToken string) (limit, offset int32) {
	// Default limit is defined as maximum page size
	limit = maxPageSize

	if pageSize > 0 && pageSize <= maxPageSize {
		limit = pageSize
	}

	if len(pageToken) > 0 {
		data, err := base64.StdEncoding.DecodeString(pageToken)
		if err != nil {
			es.log.Error("error decode offset from hash", zap.Error(err))

			return limit, offset
		}

		number, err := strconv.ParseInt(string(data), numbersBase, bitSize)
		if err != nil {
			es.log.Error("error parse offset from string", zap.Error(err))

			return limit, offset
		}

		offset = int32(number)
	}

	return limit, offset
}

func (es eventsServer) offsetPageToken(limit, offset int32) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(int(offset + limit))))
}
