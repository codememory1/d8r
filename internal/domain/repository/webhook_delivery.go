package repository

import (
	"context"

	"github.com/codememory1/d8r/internal/domain/entity"
	"github.com/codememory1/d8r/internal/domain/valueobject"
)

// WebhookDeliveryRepository defines persistence operations for webhook deliveries.
type WebhookDeliveryRepository interface {
	// GetByID returns a webhook delivery by its identifier.
	GetByID(ctx context.Context, id valueobject.ID) (*entity.WebhookDelivery, error)

	// Save persists a new webhook delivery.
	Save(ctx context.Context, delivery *entity.WebhookDelivery) error

	// Update persists changes to an existing webhook delivery.
	Update(ctx context.Context, delivery *entity.WebhookDelivery) error
}
