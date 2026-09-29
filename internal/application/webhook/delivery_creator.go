package webhook

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/codememory1/d8r/internal/application/transaction"
	"github.com/codememory1/d8r/internal/domain/entity"
	"github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/domain/valueobject"
)

// DeliveryCreator creates webhook deliveries for enabled webhooks
// subscribed to a given event.
type DeliveryCreator struct {
	tm                        transaction.Manager
	webhookRepository         repository.WebhookRepository
	webhookDeliveryRepository repository.WebhookDeliveryRepository
}

// NewDeliveryCreator creates a service for preparing webhook deliveries.
func NewDeliveryCreator(
	tm transaction.Manager,
	webhookRepository repository.WebhookRepository,
	webhookDeliveryRepository repository.WebhookDeliveryRepository,
) *DeliveryCreator {
	return &DeliveryCreator{
		tm:                        tm,
		webhookRepository:         webhookRepository,
		webhookDeliveryRepository: webhookDeliveryRepository,
	}
}

// Create saves a delivery for each matching webhook in one transaction.
func (h *DeliveryCreator) Create(ctx context.Context, eventType valueobject.WebhookEventType, payload json.RawMessage) error {
	// Select only enabled webhooks subscribed to this event type.
	webhooks, err := h.webhookRepository.FindEnabledByEventType(ctx, eventType)

	if err != nil {
		return fmt.Errorf("find subscribed webhooks: %w", err)
	}

	// Save all deliveries atomically so a failure cannot leave only
	// some subscribed webhooks with a delivery.
	return h.tm.Run(ctx, func(ctx context.Context) error {
		for _, webhook := range webhooks {
			if err != nil {
				return err
			}

			err = h.webhookDeliveryRepository.Save(ctx, entity.NewWebhookDelivery(
				webhook.ID(),
				eventType,
				payload,
			))

			if err != nil {
				return err
			}
		}

		return nil
	})
}
