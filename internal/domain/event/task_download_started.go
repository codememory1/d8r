package event

import (
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/ddd"
)

// Ensure TaskDownloadStarted implements the domain event contract.
var _ ddd.Event = (*TaskDownloadStarted)(nil)

// TaskDownloadStartedType identifies the event emitted when a task begins downloading.
const TaskDownloadStartedType ddd.EventType = "task.download.started"

// TaskDownloadStarted records when downloading begins for a task.
type TaskDownloadStarted struct {
	TaskID     valueobject.ID
	OccurredAt time.Time
}

// NewTaskDownloadStarted creates a task-download-started event.
func NewTaskDownloadStarted(taskID valueobject.ID, occurredAt time.Time) TaskDownloadStarted {
	return TaskDownloadStarted{
		TaskID:     taskID,
		OccurredAt: occurredAt,
	}
}

// Type returns the type of the task-download-started event.
func (TaskDownloadStarted) Type() ddd.EventType {
	return TaskDownloadStartedType
}
