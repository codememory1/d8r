package taskinspectionfailed

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

// CreateWebhookDeliveryHandler creates webhook deliveries for the task.inspection.failed event.
type CreateWebhookDeliveryHandler struct {
	webhookDeliveryCreator *webhook.DeliveryCreator
}

// NewCreateWebhookDeliveryHandler creates a task.inspection.failed event handler.
func NewCreateWebhookDeliveryHandler(webhookDeliveryCreator *webhook.DeliveryCreator) *CreateWebhookDeliveryHandler {
	return &CreateWebhookDeliveryHandler{
		webhookDeliveryCreator: webhookDeliveryCreator,
	}
}

// Handle creates a delivery for every enabled webhook subscribed to task.inspection.failed.
func (h *CreateWebhookDeliveryHandler) Handle(ctx context.Context, event ddd.Event) error {
	failedEvent, ok := event.(*domainevent.TaskInspectionFailed)

	if !ok {
		return fmt.Errorf("expected *event.TaskInspectionFailed, got %T", event)
	}

	webhookPayload, err := json.Marshal(payload.TaskInspectionFailed{
		TaskID: failedEvent.TaskID.String(),
	})

	if err != nil {
		return err
	}

	return h.webhookDeliveryCreator.Create(ctx, valueobject.WebhookEventTaskInspectionFailed, webhookPayload)
}
