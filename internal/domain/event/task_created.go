package event

import (
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/ddd"
)

// TaskCreatedType identifies the event emitted when a task is created.
const TaskCreatedType ddd.EventType = "task.created"

// TaskCreated represents the domain event emitted after creating a task.
type TaskCreated struct {
	ID         valueobject.ID
	OccurredAt time.Time
}

// NewTaskCreated creates a task-created domain event.
func NewTaskCreated(id valueobject.ID, occurredAt time.Time) TaskCreated {
	return TaskCreated{
		ID:         id,
		OccurredAt: occurredAt,
	}
}

// Type returns the unique type of the task-created event.
func (TaskCreated) Type() ddd.EventType {
	return TaskCreatedType
}
