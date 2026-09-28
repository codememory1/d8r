package event

import (
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/ddd"
)

// Ensure TaskInspectionStarted implements the domain event contract.
var _ ddd.Event = (*TaskInspectionStarted)(nil)

// TaskInspectionStartedType identifies the event emitted when task inspection starts.
const TaskInspectionStartedType ddd.EventType = "task.inspection.started"

// TaskInspectionStarted represents the domain event emitted when inspection
// of a task resource starts.
type TaskInspectionStarted struct {
	TaskID       valueobject.ID
	InspectionID valueobject.ID
	OccurredAt   time.Time
}

// NewTaskInspectionStarted creates a task-inspection-started domain event.
func NewTaskInspectionStarted(
	taskID valueobject.ID,
	inspectionID valueobject.ID,
	OccurredAt time.Time,
) TaskInspectionStarted {
	return TaskInspectionStarted{
		TaskID:       taskID,
		InspectionID: inspectionID,
		OccurredAt:   OccurredAt,
	}
}

// Type returns the unique type of the task-inspection-started event.
func (TaskInspectionStarted) Type() ddd.EventType {
	return TaskInspectionStartedType
}
