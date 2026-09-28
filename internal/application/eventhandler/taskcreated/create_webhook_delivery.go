package taskcreated

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/codememory1/d8r/internal/application/transaction"
	"github.com/codememory1/d8r/internal/application/webhook/payload"
	"github.com/codememory1/d8r/internal/domain/entity"
	domainevent "github.com/codememory1/d8r/internal/domain/event"
	"github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/ddd"
)

// CreateWebhookDeliveryHandler creates webhook deliveries for the task.created event.
type CreateWebhookDeliveryHandler struct {
	webhookRepository         repository.WebhookRepository
	webhookDeliveryRepository repository.WebhookDeliveryRepository
	tm                        transaction.Manager
}

// NewCreateWebhookDeliveryHandler creates a task.created event handler.
func NewCreateWebhookDeliveryHandler(
	webhookRepository repository.WebhookRepository,
	webhookDeliveryRepository repository.WebhookDeliveryRepository,
	transactionManager transaction.Manager,
) *CreateWebhookDeliveryHandler {
	return &CreateWebhookDeliveryHandler{
		webhookRepository:         webhookRepository,
		webhookDeliveryRepository: webhookDeliveryRepository,
		tm:                        transactionManager,
	}
}

// Handle creates a delivery for every enabled webhook subscribed to task.created.
func (h *CreateWebhookDeliveryHandler) Handle(ctx context.Context, event ddd.Event) error {
	createdTask, ok := event.(*domainevent.TaskCreated)

	if !ok {
		return fmt.Errorf("expected *event.TaskCreated, got %T", event)
	}

	webhooks, err := h.webhookRepository.FindEnabledByEventType(ctx, valueobject.WebhookEventTaskCreated)

	if err != nil {
		return fmt.Errorf("find subscribed webhooks: %w", err)
	}

	webhookPayload, err := json.Marshal(payload.TaskCreatedPayload{
		TaskID: createdTask.ID.String(),
	})

	return h.tm.Run(ctx, func(ctx context.Context) error {
		for _, webhook := range webhooks {
			if err != nil {
				return err
			}

			err = h.webhookDeliveryRepository.Save(ctx, entity.NewWebhookDelivery(
				webhook.ID(),
				valueobject.WebhookEventTaskCreated,
				webhookPayload,
			))

			if err != nil {
				return err
			}
		}

		return nil
	})
}
