package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/codememory1/d8r/internal/application/transaction"
	infraoutbox "github.com/codememory1/d8r/internal/infrastructure/outbox"
)

// OutboxEventWorker processes pending outbox events and dispatches them
// to the handlers registered for their event types.
type OutboxEventWorker struct {
	tm     transaction.Manager
	logger *slog.Logger
	relay  *infraoutbox.Relay
}

// NewOutboxEventWorker creates a new outbox event worker.
func NewOutboxEventWorker(
	tm transaction.Manager,
	logger *slog.Logger,
	relay *infraoutbox.Relay,
) *OutboxEventWorker {
	return &OutboxEventWorker{
		tm:     tm,
		logger: logger,
		relay:  relay,
	}
}

func (w *OutboxEventWorker) Run(ctx context.Context) error {
	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

	for {
		if err := w.relay.ProcessBatch(ctx, 10); err != nil {
			w.logger.ErrorContext(ctx, "failed to process outbox events", slog.Any("error", err))
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
}
