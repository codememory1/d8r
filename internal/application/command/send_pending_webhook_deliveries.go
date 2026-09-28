package command

import (
	"context"
	"fmt"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/cqrs"
	"golang.org/x/sync/errgroup"
)

var _ cqrs.CommandHandler[SendPendingWebhookDeliveries, struct{}] = (*SendPendingWebhookDeliveriesHandler)(nil)

// PendingWebhookDeliveryClaimer atomically claims webhook deliveries that are
// ready for processing.
type PendingWebhookDeliveryClaimer interface {
	// ClaimPending transitions pending deliveries to the processing state and
	// returns their identifiers.
	ClaimPending(ctx context.Context, limit int, maxAttempts int) ([]valueobject.ID, error)
}

// SendPendingWebhookDeliveries requests processing of a batch of pending
// webhook deliveries.
type SendPendingWebhookDeliveries struct {
	Limit int
}

// SendPendingWebhookDeliveriesHandler claims and concurrently processes
// webhook deliveries that are ready to be sent.
type SendPendingWebhookDeliveriesHandler struct {
	webhookDeliveryClaimer     PendingWebhookDeliveryClaimer
	sendWebhookDeliveryHandler cqrs.CommandHandler[SendWebhookDelivery, struct{}]
	concurrency                int
	maxAttempts                int
}

// NewSendPendingWebhookDeliveriesHandler creates a handler for processing
// batches of pending webhook deliveries.
func NewSendPendingWebhookDeliveriesHandler(
	webhookDeliveryClaimer PendingWebhookDeliveryClaimer,
	sendWebhookDeliveryHandler cqrs.CommandHandler[SendWebhookDelivery, struct{}],
	concurrency int,
	maxAttempts int,
) *SendPendingWebhookDeliveriesHandler {
	return &SendPendingWebhookDeliveriesHandler{
		webhookDeliveryClaimer:     webhookDeliveryClaimer,
		sendWebhookDeliveryHandler: sendWebhookDeliveryHandler,
		concurrency:                concurrency,
		maxAttempts:                maxAttempts,
	}
}

// Handle claims a batch of pending webhook deliveries and processes them with
// the configured concurrency limit.
func (h *SendPendingWebhookDeliveriesHandler) Handle(ctx context.Context, cmd SendPendingWebhookDeliveries) (struct{}, error) {
	webhookDeliveryIDs, err := h.webhookDeliveryClaimer.ClaimPending(ctx, cmd.Limit, h.maxAttempts)

	if err != nil {
		return struct{}{}, err
	}

	var group errgroup.Group

	group.SetLimit(h.concurrency)

	for _, webhookDeliveryID := range webhookDeliveryIDs {
		_, err := h.sendWebhookDeliveryHandler.Handle(ctx, SendWebhookDelivery{
			WebhookDeliveryID: webhookDeliveryID,
		})

		if err != nil {
			return struct{}{}, fmt.Errorf("webhook delivery %s failed: %w", webhookDeliveryID, err)
		}
	}

	return struct{}{}, group.Wait()
}
