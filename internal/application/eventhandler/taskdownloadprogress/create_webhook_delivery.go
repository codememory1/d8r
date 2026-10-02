package taskdownloadprogress

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

// CreateWebhookDeliveryHandler creates webhook deliveries for the task.download.progress event.
type CreateWebhookDeliveryHandler struct {
	webhookDeliveryCreator *webhook.DeliveryCreator
}

// NewCreateWebhookDeliveryHandler creates a task.download.progress event handler.
func NewCreateWebhookDeliveryHandler(webhookDeliveryCreator *webhook.DeliveryCreator) *CreateWebhookDeliveryHandler {
	return &CreateWebhookDeliveryHandler{
		webhookDeliveryCreator: webhookDeliveryCreator,
	}
}

// Handle creates a delivery for every enabled webhook subscribed to task.download.failed.
func (h *CreateWebhookDeliveryHandler) Handle(ctx context.Context, event ddd.Event) error {
	progressEvent, ok := event.(*domainevent.TaskDownloadProgress)

	if !ok {
		return fmt.Errorf("expected *event.TaskDownloadProgress, got %T", event)
	}

	webhookPayload, err := json.Marshal(payload.TaskDownloadProgress{
		TaskID:          progressEvent.TaskID.String(),
		DownloadedBytes: progressEvent.DownloadedBytes,
		TotalBytes:      progressEvent.TotalBytes,
		ActiveRequests:  progressEvent.ActiveRequests,
	})

	if err != nil {
		return err
	}

	return h.webhookDeliveryCreator.Create(ctx, valueobject.WebhookEventTaskDownloadProgress, webhookPayload)
}
