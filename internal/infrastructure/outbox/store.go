package outbox

import (
	"context"
)

// Store defines persistence operations required by the outbox relay.
type Store interface {
	// ClaimPending atomically claims pending outbox messages for processing.
	ClaimPending(ctx context.Context, limit int) ([]Message, error)

	// MarkProcessed marks a claimed message as successfully processed.
	MarkProcessed(ctx context.Context, message Message) error

	// MarkFailed marks a claimed message as failed.
	MarkFailed(ctx context.Context, message Message) error
}
