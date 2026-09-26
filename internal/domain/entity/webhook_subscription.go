package entity

import (
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
)

// WebhookSubscription represents a subscription to a webhook event.
type WebhookSubscription struct {
	eventType valueobject.WebhookEventType
	createdAt time.Time
}

// NewWebhookSubscription creates a new webhook subscription.
func NewWebhookSubscription(eventType valueobject.WebhookEventType, createdAt time.Time) WebhookSubscription {
	return WebhookSubscription{
		eventType: eventType,
		createdAt: createdAt,
	}
}

// EventType returns the subscribed webhook event type.
func (w WebhookSubscription) EventType() valueobject.WebhookEventType {
	return w.eventType
}

// CreatedAt returns the time when the subscription was created.
func (w WebhookSubscription) CreatedAt() time.Time {
	return w.createdAt
}
