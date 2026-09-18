package valueobject

import (
	"errors"

	"github.com/codememory1/d8r/pkg/ddd"
)

// Compile-time check that WebhookEventType implements ValueObject.
var _ ddd.ValueObject[WebhookEventType] = WebhookEventType("")

var (
	ErrInvalidWebhookEventType = errors.New("invalid webhook event type")
)

const (
	// webhookEventTaskCreated is emitted after a task has been created.
	webhookEventTaskCreated = "task.created"

	// webhookEventTaskInspectionStarted is emitted when resource inspection starts.
	webhookEventTaskInspectionStarted = "task.inspection.started"

	// webhookEventTaskInspectionCompleted is emitted after successful inspection.
	webhookEventTaskInspectionCompleted = "task.inspection.completed"

	// webhookEventTaskInspectionFailed is emitted when resource inspection fails.
	webhookEventTaskInspectionFailed = "task.inspection.failed"

	// webhookEventTaskDownloadStarted is emitted when resource downloading starts.
	webhookEventTaskDownloadStarted = "task.download.started"

	// webhookEventTaskDownloadCompleted is emitted after successful downloading.
	webhookEventTaskDownloadCompleted = "task.download.completed"

	// webhookEventTaskDownloadFailed is emitted when resource downloading fails.
	webhookEventTaskDownloadFailed = "task.download.failed"
)

// WebhookEventType identifies an event that can trigger a webhook.
type WebhookEventType string

// NewWebhookEventType creates a validated webhook event type.
func NewWebhookEventType(eventType string) (WebhookEventType, error) {
	switch eventType {
	case webhookEventTaskCreated,
		webhookEventTaskInspectionStarted,
		webhookEventTaskInspectionCompleted,
		webhookEventTaskInspectionFailed,
		webhookEventTaskDownloadStarted,
		webhookEventTaskDownloadCompleted,
		webhookEventTaskDownloadFailed:
		return WebhookEventType(eventType), nil
	default:
		return "", ErrInvalidWebhookEventType
	}
}

// String returns the string representation of the event type.
func (v WebhookEventType) String() string {
	return string(v)
}

// Equal reports whether two webhook event types are equal.
func (v WebhookEventType) Equal(other WebhookEventType) bool {
	return v == other
}
