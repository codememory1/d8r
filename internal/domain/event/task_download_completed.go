package event

import (
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/ddd"
)

// Ensure TaskCreated implements the domain event contract.
var _ ddd.Event = (*TaskDownloadCompleted)(nil)

// TaskDownloadCompletedType identifies the task download completed event.
const TaskDownloadCompletedType ddd.EventType = "task.download.completed"

// TaskDownloadCompleted represents the successful completion of a task download.
type TaskDownloadCompleted struct {
	TaskID     valueobject.ID
	OccurredAt time.Time
}

// NewTaskDownloadCompleted creates a task-download-completed domain event.
func NewTaskDownloadCompleted(taskID valueobject.ID, occurredAt time.Time) TaskDownloadCompleted {
	return TaskDownloadCompleted{
		TaskID:     taskID,
		OccurredAt: occurredAt,
	}
}

// Type returns the event type.
func (TaskDownloadCompleted) Type() ddd.EventType {
	return TaskDownloadCompletedType
}
