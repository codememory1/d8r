package taskdownloadstarted

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

// CreateWebhookDeliveryHandler creates webhook deliveries for the task.download.started event.
type CreateWebhookDeliveryHandler struct {
	webhookDeliveryCreator *webhook.DeliveryCreator
}

// NewCreateWebhookDeliveryHandler creates a task.download.started event handler.
func NewCreateWebhookDeliveryHandler(webhookDeliveryCreator *webhook.DeliveryCreator) *CreateWebhookDeliveryHandler {
	return &CreateWebhookDeliveryHandler{
		webhookDeliveryCreator: webhookDeliveryCreator,
	}
}

// Handle creates a delivery for every enabled webhook subscribed to task.created.
func (h *CreateWebhookDeliveryHandler) Handle(ctx context.Context, event ddd.Event) error {
	startedEvent, ok := event.(*domainevent.TaskDownloadStarted)

	if !ok {
		return fmt.Errorf("expected *event.TaskDownloadStarted, got %T", event)
	}

	taskID := startedEvent.TaskID.String()
	webhookPayload, err := json.Marshal(payload.TaskDownloadStarted{
		TaskID: taskID,
	})

	if err != nil {
		return err
	}

	return h.webhookDeliveryCreator.Create(
		ctx,
		valueobject.WebhookEventTaskDownloadStarted,
		&taskID,
		webhookPayload,
	)
}
