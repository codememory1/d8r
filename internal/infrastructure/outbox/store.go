package outbox

import (
	"context"
)

type Store interface {
	ClaimPending(ctx context.Context, limit int) ([]Message, error)

	MarkProcessed(ctx context.Context, message Message) error

	MarkFailed(ctx context.Context, message Message) error
}
