package taskdownloadcompleted

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/codememory1/d8r/internal/application/webhook"
	"github.com/codememory1/d8r/internal/application/webhook/payload"
	domainevent "github.com/codememory1/d8r/internal/domain/event"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/ddd"
)

// CreateWebhookDeliveryHandler creates webhook deliveries for the task.download.completed event.
type CreateWebhookDeliveryHandler struct {
	webhookDeliveryCreator *webhook.DeliveryCreator
}

// NewCreateWebhookDeliveryHandler creates a task.download.completed event handler.
func NewCreateWebhookDeliveryHandler(webhookDeliveryCreator *webhook.DeliveryCreator) *CreateWebhookDeliveryHandler {
	return &CreateWebhookDeliveryHandler{
		webhookDeliveryCreator: webhookDeliveryCreator,
	}
}

// Handle creates a delivery for every enabled webhook subscribed to task.download.completed.
func (h *CreateWebhookDeliveryHandler) Handle(ctx context.Context, event ddd.Event) error {
	createdEvent, ok := event.(*domainevent.TaskDownloadCompleted)

	if !ok {
		return fmt.Errorf("expected *event.TaskDownloadCompleted, got %T", event)
	}

	webhookPayload, err := json.Marshal(payload.TaskDownloadCompleted{
		TaskID: createdEvent.TaskID.String(),
	})

	if err != nil {
		return err
	}

	return h.webhookDeliveryCreator.Create(ctx, valueobject.WebhookEventTaskDownloadCompleted, webhookPayload)
}
