package worker

import (
	"context"
	"log/slog"
	"time"
)

// OutboxRelay processes pending outbox messages.
type OutboxRelay interface {
	ProcessBatch(ctx context.Context, limit int) error
}

// OutboxEventWorker processes pending outbox events and dispatches them
// to the handlers registered for their event types.
type OutboxEventWorker struct {
	logger *slog.Logger
	relay  OutboxRelay
	limit  int
}

// NewOutboxEventWorker creates a new outbox event worker.
func NewOutboxEventWorker(logger *slog.Logger, relay OutboxRelay, limit int) *OutboxEventWorker {
	return &OutboxEventWorker{
		logger: logger,
		relay:  relay,
		limit:  limit,
	}
}

func (w *OutboxEventWorker) Run(ctx context.Context) error {
	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

	for {
		if err := w.relay.ProcessBatch(ctx, w.limit); err != nil {
			w.logger.ErrorContext(ctx, "failed to process outbox events", slog.Any("error", err))
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
}
