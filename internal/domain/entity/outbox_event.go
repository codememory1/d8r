package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/ddd"
)

// OutboxEventStatus represents the current processing state of an outbox event.
type OutboxEventStatus string

const (
	// OutboxEventStatusPending indicates that the event is waiting to be processed.
	OutboxEventStatusPending OutboxEventStatus = "pending"

	// OutboxEventStatusProcessing indicates that the event is currently being processed.
	OutboxEventStatusProcessing OutboxEventStatus = "processing"

	// OutboxEventStatusProcessed indicates that the event was processed successfully.
	OutboxEventStatusProcessed OutboxEventStatus = "processed"

	// OutboxEventStatusFailed indicates that the event processing failed.
	OutboxEventStatusFailed OutboxEventStatus = "failed"
)

var (
	ErrInvalidOutboxEventStatus = errors.New("invalid outbox event status")
)

// OutboxEvent represents a domain event persisted for asynchronous processing.
type OutboxEvent struct {
	ddd.AggregateVersion

	id        valueobject.ID
	eventType ddd.EventType
	payload   json.RawMessage
	status    OutboxEventStatus
	createdAt time.Time
	updatedAt *time.Time
}

// NewOutboxEvent creates a new pending outbox event.
func NewOutboxEvent(eventType ddd.EventType, payload json.RawMessage) *OutboxEvent {
	return &OutboxEvent{
		id:        valueobject.NewID(),
		eventType: eventType,
		payload:   payload,
		status:    OutboxEventStatusPending,
		createdAt: time.Now(),
	}
}

// UnmarshalOutboxEvent restores an outbox event from persisted data.
func UnmarshalOutboxEvent(
	id valueobject.ID,
	eventType ddd.EventType,
	payload json.RawMessage,
	status string,
	version int64,
	createdAt time.Time,
	updatedAt *time.Time,
) (*OutboxEvent, error) {
	outboxEventStatus, err := ParseOutboxEventStatus(status)

	if err != nil {
		return nil, err
	}

	return &OutboxEvent{
		AggregateVersion: ddd.NewAggregateVersion(version),
		id:               id,
		eventType:        eventType,
		payload:          payload,
		status:           outboxEventStatus,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}, nil
}

// ID returns the unique identifier of the outbox event.
func (e *OutboxEvent) ID() valueobject.ID {
	return e.id
}

// EventType returns the type of the persisted domain event.
func (e *OutboxEvent) EventType() ddd.EventType {
	return e.eventType
}

// Payload returns the serialized event payload.
func (e *OutboxEvent) Payload() json.RawMessage {
	return e.payload
}

// Status returns the current processing status of the outbox event.
func (e *OutboxEvent) Status() OutboxEventStatus {
	return e.status
}

// CreatedAt returns the time at which the outbox event was created.
func (e *OutboxEvent) CreatedAt() time.Time {
	return e.createdAt
}

// UpdatedAt returns the time the event was updated in the outbox.
func (e *OutboxEvent) UpdatedAt() *time.Time {
	return e.updatedAt
}

// Processing marks the outbox event as being processed.
func (e *OutboxEvent) Processing() {
	e.status = OutboxEventStatusProcessing
}

// Processed marks the outbox event as successfully processed.
func (e *OutboxEvent) Processed() {
	e.status = OutboxEventStatusProcessed
}

// Failed marks the outbox event as failed.
func (e *OutboxEvent) Failed() {
	e.status = OutboxEventStatusFailed
}

func ParseOutboxEventStatus(
	value string,
) (OutboxEventStatus, error) {
	status := OutboxEventStatus(value)

	switch status {
	case OutboxEventStatusPending,
		OutboxEventStatusProcessing,
		OutboxEventStatusProcessed,
		OutboxEventStatusFailed:
		return status, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidOutboxEventStatus, value)
	}
}
