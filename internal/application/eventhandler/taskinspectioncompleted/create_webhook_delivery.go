package taskinspectioncompleted

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

// CreateWebhookDeliveryHandler creates webhook deliveries for the task.inspection.completed event.
type CreateWebhookDeliveryHandler struct {
	webhookDeliveryCreator *webhook.DeliveryCreator
}

// NewCreateWebhookDeliveryHandler creates a task.created event handler.
func NewCreateWebhookDeliveryHandler(webhookDeliveryCreator *webhook.DeliveryCreator) *CreateWebhookDeliveryHandler {
	return &CreateWebhookDeliveryHandler{
		webhookDeliveryCreator: webhookDeliveryCreator,
	}
}

// Handle creates a delivery for every enabled webhook subscribed to task.inspection.completed.
func (h *CreateWebhookDeliveryHandler) Handle(ctx context.Context, event ddd.Event) error {
	completedEvent, ok := event.(*domainevent.TaskInspectionCompleted)

	if !ok {
		return fmt.Errorf("expected *event.TaskInspectionCompleted, got %T", event)
	}

	taskID := completedEvent.InspectionID.String()
	webhookPayload, err := json.Marshal(payload.TaskInspectionCompleted{
		TaskID:       taskID,
		InspectionID: completedEvent.InspectionID.String(),
	})

	if err != nil {
		return err
	}

	return h.webhookDeliveryCreator.Create(
		ctx,
		valueobject.WebhookEventTaskInspectionCompleted,
		&taskID,
		webhookPayload,
	)
}
