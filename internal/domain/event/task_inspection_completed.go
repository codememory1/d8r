package event

import (
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/ddd"
)

// TaskInspectionCompletedType identifies the event emitted after a task
// inspection completes successfully.
const TaskInspectionCompletedType ddd.EventType = "task.inspection.completed"

// TaskInspectionCompleted represents a successfully completed task inspection.
type TaskInspectionCompleted struct {
	TaskID       valueobject.ID
	InspectionID valueobject.ID
	OccurredAt   time.Time
}

// NewTaskInspectionCompleted creates a task-inspection-completed event.
func NewTaskInspectionCompleted(taskID valueobject.ID, inspectionID valueobject.ID, occurredAt time.Time) TaskInspectionCompleted {
	return TaskInspectionCompleted{
		TaskID:       taskID,
		InspectionID: inspectionID,
		OccurredAt:   occurredAt,
	}
}

// Type returns the type of the task-inspection-completed event.
func (TaskInspectionCompleted) Type() ddd.EventType {
	return TaskInspectionCompletedType
}
