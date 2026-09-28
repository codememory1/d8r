package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/ddd"
	"github.com/codememory1/d8r/pkg/statemachine"
)

var _ ddd.Entity[valueobject.ID] = (*WebhookDelivery)(nil)

// WebhookDeliveryTransition identifies a state transition of a webhook delivery.
type WebhookDeliveryTransition string

// WebhookDeliveryStatus represents the current state of a webhook delivery.
type WebhookDeliveryStatus string

const (
	// WebhookDeliveryTransitionProcess starts processing the delivery.
	WebhookDeliveryTransitionProcess WebhookDeliveryTransition = "process"

	// WebhookDeliveryTransitionRetry schedules the delivery for another attempt.
	WebhookDeliveryTransitionRetry WebhookDeliveryTransition = "retry"

	// WebhookDeliveryTransitionDeliver marks the delivery as successfully delivered.
	WebhookDeliveryTransitionDeliver WebhookDeliveryTransition = "deliver"

	// WebhookDeliveryTransitionFail marks the delivery as permanently failed.
	WebhookDeliveryTransitionFail WebhookDeliveryTransition = "fail"
)

const (
	// WebhookDeliveryStatusPending indicates that the delivery is waiting to be processed.
	WebhookDeliveryStatusPending WebhookDeliveryStatus = "pending"

	// WebhookDeliveryStatusProcessing indicates that the delivery is being processed.
	WebhookDeliveryStatusProcessing WebhookDeliveryStatus = "processing"

	// WebhookDeliveryStatusDelivered indicates that the delivery was completed successfully.
	WebhookDeliveryStatusDelivered WebhookDeliveryStatus = "delivered"

	// WebhookDeliveryStatusFailed indicates that the delivery permanently failed.
	WebhookDeliveryStatusFailed WebhookDeliveryStatus = "failed"
)

var webhookDeliveryStateMachine = statemachine.NewMachine(
	statemachine.T(
		WebhookDeliveryTransitionProcess,
		[]WebhookDeliveryStatus{
			WebhookDeliveryStatusPending,
		},
		WebhookDeliveryStatusProcessing,
	),
	statemachine.T(
		WebhookDeliveryTransitionRetry,
		[]WebhookDeliveryStatus{
			WebhookDeliveryStatusProcessing,
		},
		WebhookDeliveryStatusPending,
	),
	statemachine.T(
		WebhookDeliveryTransitionDeliver,
		[]WebhookDeliveryStatus{
			WebhookDeliveryStatusProcessing,
		},
		WebhookDeliveryStatusDelivered,
	),
	statemachine.T(
		WebhookDeliveryTransitionFail,
		[]WebhookDeliveryStatus{
			WebhookDeliveryStatusProcessing,
		},
		WebhookDeliveryStatusFailed,
	),
)

var (
	// ErrInvalidWebhookDeliveryStatus is returned when a webhook delivery status cannot be recognized.
	ErrInvalidWebhookDeliveryStatus = errors.New("invalid webhook delivery status")
)

// WebhookDelivery represents an attempt to deliver a webhook event to a subscribed webhook.
type WebhookDelivery struct {
	ddd.AggregateVersion

	id             valueobject.ID
	webhookID      valueobject.ID
	eventType      valueobject.WebhookEventType
	payload        json.RawMessage
	status         WebhookDeliveryStatus
	attempts       int
	nextAttemptAt  time.Time
	responseStatus *int
	lastError      *string
	createdAt      time.Time
	updatedAt      *time.Time
	deliveredAt    *time.Time
}

// NewWebhookDelivery creates a pending webhook delivery ready for its first attempt.
func NewWebhookDelivery(
	webhookID valueobject.ID,
	eventType valueobject.WebhookEventType,
	payload json.RawMessage,
) *WebhookDelivery {
	now := time.Now()

	return &WebhookDelivery{
		id:            valueobject.NewID(),
		webhookID:     webhookID,
		eventType:     eventType,
		payload:       append(json.RawMessage(nil), payload...),
		status:        WebhookDeliveryStatusPending,
		attempts:      0,
		nextAttemptAt: now,
		createdAt:     now,
	}
}

// UnmarshalWebhookDelivery reconstructs a webhook delivery entity from its
// persisted state.
func UnmarshalWebhookDelivery(
	id valueobject.ID,
	webhookID valueobject.ID,
	eventType valueobject.WebhookEventType,
	payload json.RawMessage,
	status string,
	attempts int,
	nextAttemptAt time.Time,
	responseStatus *int,
	lastError *string,
	version int64,
	createdAt time.Time,
	updatedAt *time.Time,
	deliveredAt *time.Time,
) (*WebhookDelivery, error) {
	webhookDeliveryStatus, err := ParseWebhookDeliveryStatus(status)

	if err != nil {
		return nil, err
	}

	return &WebhookDelivery{
		AggregateVersion: ddd.NewAggregateVersion(version),
		id:               id,
		webhookID:        webhookID,
		eventType:        eventType,
		payload:          payload,
		status:           webhookDeliveryStatus,
		attempts:         attempts,
		nextAttemptAt:    nextAttemptAt,
		responseStatus:   responseStatus,
		lastError:        lastError,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
		deliveredAt:      deliveredAt,
	}, nil
}

// ID returns the delivery identifier.
func (w *WebhookDelivery) ID() valueobject.ID {
	return w.id
}

// WebhookID returns the identifier of the target webhook.
func (w *WebhookDelivery) WebhookID() valueobject.ID {
	return w.webhookID
}

// EventType returns the type of webhook event being delivered.
func (w *WebhookDelivery) EventType() valueobject.WebhookEventType {
	return w.eventType
}

// Payload returns a copy of the serialized webhook event payload.
func (w *WebhookDelivery) Payload() json.RawMessage {
	return append(json.RawMessage(nil), w.payload...)
}

// Status returns the current delivery status.
func (w *WebhookDelivery) Status() WebhookDeliveryStatus {
	return w.status
}

// Attempts returns the number of completed delivery attempts.
func (w *WebhookDelivery) Attempts() int {
	return w.attempts
}

// NextAttemptAt returns the earliest time at which another attempt may be made.
func (w *WebhookDelivery) NextAttemptAt() time.Time {
	return w.nextAttemptAt
}

// ResponseStatus returns a copy of the HTTP status received during the latest
// delivery attempt, or nil when no HTTP response was received.
func (w *WebhookDelivery) ResponseStatus() *int {
	if w.responseStatus == nil {
		return nil
	}

	return new(*w.responseStatus)
}

// LastError returns a copy of the error message from the latest failed attempt.
func (w *WebhookDelivery) LastError() *string {
	if w.lastError == nil {
		return nil
	}

	return new(*w.lastError)
}

// CreatedAt returns the time at which the delivery was created.
func (w *WebhookDelivery) CreatedAt() time.Time {
	return w.createdAt
}

// UpdatedAt returns a copy of the time at which the delivery was last updated.
func (w *WebhookDelivery) UpdatedAt() *time.Time {
	if w.updatedAt == nil {
		return nil
	}

	return new(*w.updatedAt)
}

// DeliveredAt returns a copy of the time at which the delivery succeeded.
func (w *WebhookDelivery) DeliveredAt() *time.Time {
	if w.deliveredAt == nil {
		return nil
	}

	return new(*w.deliveredAt)
}

// Process transitions the pending delivery to the processing state.
func (w *WebhookDelivery) Process() error {
	return w.transition(WebhookDeliveryTransitionProcess)
}

// Retry records an unsuccessful attempt and returns the delivery to the pending state.
func (w *WebhookDelivery) Retry(nextAttemptAt time.Time, responseStatus *int, lastError string) error {
	err := w.transition(WebhookDeliveryTransitionRetry)

	if err != nil {
		return err
	}

	w.attempts++
	w.nextAttemptAt = nextAttemptAt
	w.responseStatus = responseStatus
	w.lastError = &lastError

	return nil
}

// Deliver records a successful attempt and marks the delivery as delivered.
func (w *WebhookDelivery) Deliver(responseStatus *int) error {
	err := w.transition(WebhookDeliveryTransitionDeliver)

	if err != nil {
		return err
	}

	w.attempts++
	w.responseStatus = responseStatus
	w.lastError = nil
	w.deliveredAt = new(time.Now())

	return nil
}

// Fail records an unsuccessful final attempt and marks the delivery as permanently failed.
func (w *WebhookDelivery) Fail(responseStatus *int, lastError string) error {
	err := w.transition(WebhookDeliveryTransitionFail)

	if err != nil {
		return err
	}

	w.attempts++
	w.responseStatus = responseStatus
	w.lastError = &lastError

	return nil
}

// transition applies the named state transition to the webhook delivery.
func (w *WebhookDelivery) transition(name WebhookDeliveryTransition) error {
	newStatus, err := webhookDeliveryStateMachine.Transition(name, w.status)

	if err != nil {
		return err
	}

	w.status = newStatus
	w.updatedAt = new(time.Now())

	return nil
}

// ParseWebhookDeliveryStatus parses and validates a webhook delivery status.
func ParseWebhookDeliveryStatus(value string) (WebhookDeliveryStatus, error) {
	status := WebhookDeliveryStatus(value)

	switch status {
	case WebhookDeliveryStatusPending,
		WebhookDeliveryStatusProcessing,
		WebhookDeliveryStatusDelivered,
		WebhookDeliveryStatusFailed:
		return status, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidWebhookDeliveryStatus, value)
	}
}
