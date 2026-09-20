package worker

import (
	"context"
	"log/slog"
	"time"

	appevent "github.com/codememory1/d8r/internal/application/event"
	"github.com/codememory1/d8r/internal/application/transaction"
	"github.com/codememory1/d8r/internal/domain/entity"
	"github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/infrastructure/eventcodec"
	"golang.org/x/sync/errgroup"
)

// OutboxEventWorker processes pending outbox events and dispatches them
// to the handlers registered for their event types.
type OutboxEventWorker struct {
	tm                    transaction.Manager
	logger                *slog.Logger
	eventDecoder          *eventcodec.Decoder
	eventDispatcher       appevent.Dispatcher
	outboxEventRepository repository.OutboxEventRepository
}

// NewOutboxEventWorker creates a new outbox event worker.
func NewOutboxEventWorker(
	tm transaction.Manager,
	logger *slog.Logger,
	eventDecoder *eventcodec.Decoder,
	eventDispatcher appevent.Dispatcher,
	outboxEventRepository repository.OutboxEventRepository,
) *OutboxEventWorker {
	return &OutboxEventWorker{
		tm:                    tm,
		logger:                logger,
		eventDecoder:          eventDecoder,
		eventDispatcher:       eventDispatcher,
		outboxEventRepository: outboxEventRepository,
	}
}

// Run continuously claims and processes pending outbox events until
// the context is canceled or an unrecoverable error occurs.
func (w *OutboxEventWorker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			outboxEvents, err := w.outboxEventRepository.ClaimPending(ctx, 10)

			if err != nil {
				return err
			}

			if len(outboxEvents) == 0 {
				if err := w.wait(ctx); err != nil {
					return err
				}

				continue
			}

			if err := w.processOutboxEvents(ctx, outboxEvents); err != nil {
				return err
			}
		}
	}
}

// processOutboxEvents processes a batch of outbox events concurrently.
func (w *OutboxEventWorker) processOutboxEvents(ctx context.Context, outboxEvents []*entity.OutboxEvent) error {
	var group errgroup.Group

	group.SetLimit(10)

	for _, outboxEvent := range outboxEvents {
		group.Go(func() error {
			return w.processOutboxEvent(ctx, outboxEvent)
		})
	}

	return group.Wait()
}

// processOutboxEvent decodes, dispatches, and updates the state of a single outbox event.
func (w *OutboxEventWorker) processOutboxEvent(ctx context.Context, outboxEvent *entity.OutboxEvent) error {
	// Restore the original domain event from the persisted outbox payload.
	decodedEvent, err := w.eventDecoder.Decode(outboxEvent.EventType(), outboxEvent.Payload())

	if err != nil {
		w.logger.ErrorContext(
			ctx,
			"failed to decode outbox event",
			slog.String("outbox_event_id", outboxEvent.ID().String()),
			slog.Any("error", err),
		)

		return nil
	}

	// Dispatch the event to all registered handlers.
	if dispatchErr := w.eventDispatcher.Dispatch(ctx, *decodedEvent); dispatchErr != nil {
		w.logger.ErrorContext(
			ctx,
			"failed to dispatch outbox event",
			slog.String("outbox_event_id", outboxEvent.ID().String()),
			slog.Any("error", dispatchErr),
		)

		outboxEvent.Failed()

		if updateErr := w.outboxEventRepository.Update(ctx, outboxEvent); updateErr != nil {
			w.logger.ErrorContext(
				ctx,
				"failed to update outbox event to failed state",
				slog.String("outbox_event_id", outboxEvent.ID().String()),
				slog.Any("error", updateErr),
			)
		}

		return nil
	}

	// Mark the event as processed only after all handlers complete successfully.
	outboxEvent.Processed()

	if updateErr := w.outboxEventRepository.Update(ctx, outboxEvent); updateErr != nil {
		w.logger.ErrorContext(
			ctx,
			"failed to update outbox event to processed state",
			slog.String("outbox_event_id", outboxEvent.ID().String()),
			slog.Any("error", updateErr),
		)
	}

	return nil
}

// wait pauses polling until the retry interval expires or the context is canceled.
func (w *OutboxEventWorker) wait(ctx context.Context) error {
	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
