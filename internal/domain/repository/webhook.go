package repository

import (
	"context"

	"github.com/codememory1/d8r/internal/domain/entity"
	"github.com/codememory1/d8r/internal/domain/valueobject"
)

// WebhookRepository defines persistence operations for webhook aggregates.
type WebhookRepository interface {
	// GetByID returns a webhook aggregate by its identifier.
	GetByID(ctx context.Context, id valueobject.ID) (*entity.Webhook, error)

	// FindEnabledByEventType returns all enabled webhooks subscribed to the specified event type.
	FindEnabledByEventType(ctx context.Context, eventType valueobject.WebhookEventType) ([]*entity.Webhook, error)

	// Save persists a webhook and its subscriptions.
	Save(ctx context.Context, webhook *entity.Webhook) error
}
