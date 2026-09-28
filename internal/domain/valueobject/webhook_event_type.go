package valueobject

import (
	"errors"

	"github.com/codememory1/d8r/pkg/ddd"
)

// Compile-time check that WebhookEventType implements ValueObject.
var _ ddd.ValueObject[WebhookEventType] = WebhookEventType("")

var (
	// ErrInvalidWebhookEventType is returned when an unsupported webhook event
	// type is provided.
	ErrInvalidWebhookEventType = errors.New("invalid webhook event type")
)

const (
	// WebhookEventTaskCreated is emitted after a task has been created.
	WebhookEventTaskCreated = "task.created"

	// WebhookEventTaskInspectionStarted is emitted when resource inspection starts.
	WebhookEventTaskInspectionStarted = "task.inspection.started"

	// WebhookEventTaskInspectionCompleted is emitted after successful inspection.
	WebhookEventTaskInspectionCompleted = "task.inspection.completed"

	// WebhookEventTaskInspectionFailed is emitted when resource inspection fails.
	WebhookEventTaskInspectionFailed = "task.inspection.failed"

	// WebhookEventTaskDownloadStarted is emitted when resource downloading starts.
	WebhookEventTaskDownloadStarted = "task.download.started"

	// WebhookEventTaskDownloadCompleted is emitted after successful downloading.
	WebhookEventTaskDownloadCompleted = "task.download.completed"

	// WebhookEventTaskDownloadFailed is emitted when resource downloading fails.
	WebhookEventTaskDownloadFailed = "task.download.failed"
)

// WebhookEventType identifies an event that can trigger a webhook.
type WebhookEventType string

// NewWebhookEventType creates a validated webhook event type.
func NewWebhookEventType(eventType string) (WebhookEventType, error) {
	switch eventType {
	case WebhookEventTaskCreated,
		WebhookEventTaskInspectionStarted,
		WebhookEventTaskInspectionCompleted,
		WebhookEventTaskInspectionFailed,
		WebhookEventTaskDownloadStarted,
		WebhookEventTaskDownloadCompleted,
		WebhookEventTaskDownloadFailed:
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
