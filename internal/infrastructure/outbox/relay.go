package outbox

import (
	"context"

	appevent "github.com/codememory1/d8r/internal/application/event"
	"golang.org/x/sync/errgroup"
)

// Relay claims outbox messages, reconstructs their domain events and dispatches
// them to registered handlers.
type Relay struct {
	registry    EventRegistry
	decoder     Decoder
	store       Store
	dispatcher  appevent.Dispatcher
	concurrency int
}

// NewRelay creates an outbox relay with bounded concurrent processing.
func NewRelay(
	registry EventRegistry,
	decoder Decoder,
	store Store,
	dispatcher appevent.Dispatcher,
	concurrency int,
) *Relay {
	return &Relay{
		registry:    registry,
		decoder:     decoder,
		store:       store,
		dispatcher:  dispatcher,
		concurrency: concurrency,
	}
}

// ProcessBatch claims and processes up to the specified number of pending
// outbox messages.
func (r *Relay) ProcessBatch(ctx context.Context, limit int) error {
	messages, err := r.store.ClaimPending(ctx, limit)

	if err != nil {
		return err
	}

	var group errgroup.Group

	group.SetLimit(r.concurrency)

	for _, message := range messages {
		group.Go(func() error {
			return r.processMessage(ctx, message)
		})
	}

	return group.Wait()
}

// processMessage reconstructs and dispatches one outbox event, then persists
// its final processing status.
func (r *Relay) processMessage(ctx context.Context, message Message) error {
	event, err := r.registry.Get(message.EventType)

	if err != nil {
		return r.store.MarkFailed(ctx, message)
	}

	if decodeErr := r.decoder.Decode(message.Payload, event); decodeErr != nil {
		return r.store.MarkFailed(ctx, message)
	}

	if err := r.dispatcher.Dispatch(ctx, event); err != nil {
		return r.store.MarkFailed(ctx, message)
	}

	if markProcessedErr := r.store.MarkProcessed(ctx, message); markProcessedErr != nil {
		return markProcessedErr
	}

	return nil
}
