package entity

import (
	"errors"
	"fmt"
	"maps"
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/statemachine"
)

// WebhookTransition represents a named webhook state transition.
type WebhookTransition string

// WebhookStatus represents the operational state of a webhook.
type WebhookStatus string

const (
	// WebhookTransitionEnable enable webhook.
	WebhookTransitionEnable WebhookTransition = "enable"

	// WebhookTransitionDisable disable webhook.
	WebhookTransitionDisable WebhookTransition = "disable"
)

const (
	// WebhookStatusEnabled indicates that the webhook can receive events.
	WebhookStatusEnabled WebhookStatus = "enabled"

	// WebhookStatusDisabled indicates that event delivery is disabled.
	WebhookStatusDisabled WebhookStatus = "disabled"
)

// webhookStateMachine defines allowed webhook status transitions.
var webhookStateMachine = statemachine.NewMachine(
	statemachine.T(
		WebhookTransitionEnable,
		[]WebhookStatus{
			WebhookStatusDisabled,
		},
		WebhookStatusEnabled,
	),
	statemachine.T(
		WebhookTransitionDisable,
		[]WebhookStatus{
			WebhookStatusEnabled,
		},
		WebhookStatusDisabled,
	),
)

var (
	// ErrWebhookAlreadySubscribed is returned when the webhook is already
	// subscribed to the requested event type.
	ErrWebhookAlreadySubscribed = errors.New("webhook already subscribed")

	// ErrWebhookUnsubscribed is returned when the requested subscription
	// does not exist.
	ErrWebhookUnsubscribed = errors.New("webhook subscription not found")

	ErrInvalidWebhookStatus = errors.New("invalid webhook status")
)

// Webhook is an aggregate root that represents an HTTP endpoint and its
// subscriptions to application events.
type Webhook struct {
	id            valueobject.ID
	url           valueobject.URL
	headers       valueobject.Headers
	status        WebhookStatus
	subscriptions map[valueobject.WebhookEventType]WebhookSubscription
	createdAt     time.Time
	updatedAt     *time.Time
}

// NewWebhook creates an enabled webhook with the provided subscriptions.
func NewWebhook(
	url valueobject.URL,
	headers valueobject.Headers,
	subscriptions []valueobject.WebhookEventType,
) (*Webhook, error) {
	webhook := &Webhook{
		id:            valueobject.NewID(),
		url:           url,
		headers:       headers,
		status:        WebhookStatusEnabled,
		subscriptions: make(map[valueobject.WebhookEventType]WebhookSubscription, len(subscriptions)),
		createdAt:     time.Now(),
	}

	for _, sub := range subscriptions {
		if err := webhook.Subscribe(sub); err != nil {
			return nil, err
		}
	}

	return webhook, nil
}

// UnmarshalWebhook reconstructs a webhook aggregate from persisted data.
func UnmarshalWebhook(
	id valueobject.ID,
	url valueobject.URL,
	headers valueobject.Headers,
	status string,
	subscriptions map[valueobject.WebhookEventType]WebhookSubscription,
	createdAt time.Time,
	updatedAt *time.Time,
) (*Webhook, error) {
	webhookStatus, err := ParseWebhookStatus(status)

	if err != nil {
		return nil, err
	}

	clonedSubscriptions := maps.Clone(subscriptions)

	if clonedSubscriptions == nil {
		clonedSubscriptions = make(map[valueobject.WebhookEventType]WebhookSubscription)
	}

	return &Webhook{
		id:            id,
		url:           url,
		headers:       headers,
		status:        webhookStatus,
		subscriptions: clonedSubscriptions,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}, nil
}

// ID returns the webhook identifier.
func (w *Webhook) ID() valueobject.ID {
	return w.id
}

// URL returns the webhook endpoint URL.
func (w *Webhook) URL() valueobject.URL {
	return w.url
}

// Headers returns the HTTP headers configured for webhook delivery.
func (w *Webhook) Headers() valueobject.Headers {
	return w.headers
}

// Status returns the current webhook status.
func (w *Webhook) Status() WebhookStatus {
	return w.status
}

// Subscriptions returns the webhook's event subscriptions.
func (w *Webhook) Subscriptions() map[valueobject.WebhookEventType]WebhookSubscription {
	return maps.Clone(w.subscriptions)
}

// CreatedAt returns the time when the webhook was created.
func (w *Webhook) CreatedAt() time.Time {
	return w.createdAt
}

// UpdatedAt returns the time when the webhook was last updated.
func (w *Webhook) UpdatedAt() *time.Time {
	return w.updatedAt
}

// Subscribe subscribes the webhook to the specified event type.
func (w *Webhook) Subscribe(eventType valueobject.WebhookEventType) error {
	if _, ok := w.subscriptions[eventType]; ok {
		return ErrWebhookAlreadySubscribed
	}

	w.subscriptions[eventType] = NewWebhookSubscription(eventType, time.Now())

	return nil
}

// Unsubscribe removes the webhook subscription for the specified event type.
func (w *Webhook) Unsubscribe(eventType valueobject.WebhookEventType) error {
	if _, ok := w.subscriptions[eventType]; !ok {
		return ErrWebhookUnsubscribed
	}

	delete(w.subscriptions, eventType)

	return nil
}

// Supports reports whether the webhook is subscribed to the specified event type.
func (w *Webhook) Supports(eventType valueobject.WebhookEventType) bool {
	_, ok := w.subscriptions[eventType]

	return ok
}

// Enable enables event delivery to the webhook.
func (w *Webhook) Enable() error {
	return w.transition(WebhookTransitionEnable)
}

// Disable disables event delivery to the webhook.
func (w *Webhook) Disable() error {
	return w.transition(WebhookTransitionDisable)
}

// transition applies the named state transition to the webhook.
func (w *Webhook) transition(name WebhookTransition) error {
	newStatus, err := webhookStateMachine.Transition(name, w.status)

	if err != nil {
		return err
	}

	w.status = newStatus
	w.updatedAt = new(time.Now())

	return nil
}

func ParseWebhookStatus(value string) (WebhookStatus, error) {
	status := WebhookStatus(value)

	switch status {
	case WebhookStatusEnabled,
		WebhookStatusDisabled:
		return status, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidWebhookStatus, value)
	}
}
