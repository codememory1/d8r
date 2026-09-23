package outbox

import (
	"context"

	appevent "github.com/codememory1/d8r/internal/application/event"
	"golang.org/x/sync/errgroup"
)

type Relay struct {
	store       Store
	decoder     Decoder
	dispatcher  appevent.Dispatcher
	concurrency int
}

func NewRelay(
	store Store,
	decoder Decoder,
	dispatcher appevent.Dispatcher,
	concurrency int,
) *Relay {
	return &Relay{
		store:       store,
		decoder:     decoder,
		dispatcher:  dispatcher,
		concurrency: concurrency,
	}
}

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

func (r *Relay) processMessage(ctx context.Context, message Message) error {
	event, err := r.decoder.Decode(message.EventType, message.Payload)

	if err != nil {
		if markFailedErr := r.store.MarkFailed(ctx, message); markFailedErr != nil {
			return markFailedErr
		}

		return nil
	}

	if err := r.dispatcher.Dispatch(ctx, event); err != nil {
		if markFailedErr := r.store.MarkFailed(ctx, message); markFailedErr != nil {
			return markFailedErr
		}

		return nil
	}

	if markProcessedErr := r.store.MarkProcessed(ctx, message); markProcessedErr != nil {
		return markProcessedErr
	}

	return nil
}
