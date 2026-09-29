package taskinspectionstarted

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

// CreateWebhookDeliveryHandler creates webhook deliveries for the task.inspection.started event.
type CreateWebhookDeliveryHandler struct {
	webhookDeliveryCreator *webhook.DeliveryCreator
}

// NewCreateWebhookDeliveryHandler creates a task.inspection.started event handler.
func NewCreateWebhookDeliveryHandler(webhookDeliveryCreator *webhook.DeliveryCreator) *CreateWebhookDeliveryHandler {
	return &CreateWebhookDeliveryHandler{
		webhookDeliveryCreator: webhookDeliveryCreator,
	}
}

// Handle creates a delivery for every enabled webhook subscribed to task.created.
func (h *CreateWebhookDeliveryHandler) Handle(ctx context.Context, event ddd.Event) error {
	startedEvent, ok := event.(*domainevent.TaskInspectionStarted)

	if !ok {
		return fmt.Errorf("expected *event.TaskInspectionStarted, got %T", event)
	}

	webhookPayload, err := json.Marshal(payload.TaskInspectionStarted{
		TaskID: startedEvent.TaskID.String(),
	})

	if err != nil {
		return err
	}

	return h.webhookDeliveryCreator.Create(ctx, valueobject.WebhookEventTaskInspectionStarted, webhookPayload)
}
