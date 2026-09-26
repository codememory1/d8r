package repository

import (
	"context"

	"github.com/codememory1/d8r/internal/domain/entity"
)

// WebhookRepository defines persistence operations for webhook aggregates.
type WebhookRepository interface {
	// Save persists a webhook and its subscriptions.
	Save(ctx context.Context, webhook *entity.Webhook) error
}
