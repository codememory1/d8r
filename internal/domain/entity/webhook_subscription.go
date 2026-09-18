package entity

import (
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
)

type WebhookSubscription struct {
	eventType valueobject.WebhookEventType
	createdAt time.Time
}

func NewWebhookSubscription(eventType valueobject.WebhookEventType, createdAt time.Time) WebhookSubscription {
	return WebhookSubscription{
		eventType: eventType,
		createdAt: createdAt,
	}
}

func (w WebhookSubscription) EventType() valueobject.WebhookEventType {
	return w.eventType
}

func (w WebhookSubscription) CreatedAt() time.Time {
	return w.createdAt
}
