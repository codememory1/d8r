package repository

import (
	"context"

	"github.com/codememory1/d8r/internal/domain/entity"
)

// OutboxEventRepository defines persistence operations for outbox events.
type OutboxEventRepository interface {
	// ClaimPending atomically claims pending outbox events for processing.
	ClaimPending(ctx context.Context, limit int) ([]*entity.OutboxEvent, error)

	// Update persists changes to an existing outbox event.
	Update(ctx context.Context, outboxEvent *entity.OutboxEvent) error

	// Save persists a new outbox event.
	Save(ctx context.Context, outboxEvent *entity.OutboxEvent) error
}
