package event

import (
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/ddd"
)

var _ ddd.Event = (*TaskInspectionFailed)(nil)

const TaskInspectionFailedType ddd.EventType = "task.inspection.failed"

type TaskInspectionFailed struct {
	TaskID     valueobject.ID
	OccurredAt time.Time
}

func NewTaskInspectionFailed(taskID valueobject.ID, occurredAt time.Time) TaskInspectionFailed {
	return TaskInspectionFailed{
		TaskID:     taskID,
		OccurredAt: occurredAt,
	}
}

func (TaskInspectionFailed) Type() ddd.EventType {
	return TaskInspectionFailedType
}
