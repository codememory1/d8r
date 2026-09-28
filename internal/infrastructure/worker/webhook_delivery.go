package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/codememory1/d8r/internal/application/command"
	"github.com/codememory1/d8r/pkg/cqrs"
)

// WebhookDeliveryWorker periodically processes webhook deliveries that are
// ready to be sent.
type WebhookDeliveryWorker struct {
	logger                              *slog.Logger
	sendPendingWebhookDeliveriesHandler cqrs.CommandHandler[command.SendPendingWebhookDeliveries, struct{}]
	limit                               int
}

// NewWebhookDeliveryWorker creates a worker for processing pending webhook
// deliveries.
func NewWebhookDeliveryWorker(
	logger *slog.Logger,
	sendPendingWebhookDeliveriesHandler cqrs.CommandHandler[command.SendPendingWebhookDeliveries, struct{}],
	limit int,
) *WebhookDeliveryWorker {
	return &WebhookDeliveryWorker{
		logger:                              logger,
		sendPendingWebhookDeliveriesHandler: sendPendingWebhookDeliveriesHandler,
		limit:                               limit,
	}
}

// Run continuously processes pending webhook deliveries until the context is
// canceled.
func (w *WebhookDeliveryWorker) Run(ctx context.Context) error {
	for {
		_, err := w.sendPendingWebhookDeliveriesHandler.Handle(ctx, command.SendPendingWebhookDeliveries{
			Limit: w.limit,
		})

		if err != nil {
			w.logger.ErrorContext(
				ctx,
				"Failed to send pending webhook deliveries",
				slog.Any("error", err),
			)
		}

		timer := time.NewTimer(1 * time.Second)

		select {
		case <-ctx.Done():
			timer.Stop()

			return ctx.Err()
		case <-timer.C:
		}
	}
}
