package event

import (
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/ddd"
)

// Ensure TaskDownloadCompleted implements the domain event contract.
var _ ddd.Event = (*TaskDownloadFailed)(nil)

// TaskDownloadFailedType identifies the task download failed event.
const TaskDownloadFailedType ddd.EventType = "task.download.failed"

// TaskDownloadFailed is a domain event indicating that a task download failed.
type TaskDownloadFailed struct {
	TaskID     valueobject.ID
	OccurredAt time.Time
}

// NewTaskDownloadFailed creates a task-download-failed domain event.
func NewTaskDownloadFailed(taskID valueobject.ID, occurredAt time.Time) TaskDownloadFailed {
	return TaskDownloadFailed{
		TaskID:     taskID,
		OccurredAt: occurredAt,
	}
}

// Type returns the event type.
func (TaskDownloadFailed) Type() ddd.EventType {
	return TaskDownloadFailedType
}
