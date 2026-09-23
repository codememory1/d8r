package event

import (
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/ddd"
)

const TaskCreatedType ddd.EventType = "task.created"

type TaskCreated struct {
	ID         valueobject.ID
	OccurredAt time.Time
}

func NewTaskCreated(id valueobject.ID, occurredAt time.Time) TaskCreated {
	return TaskCreated{
		ID:         id,
		OccurredAt: occurredAt,
	}
}

func (TaskCreated) Type() ddd.EventType {
	return TaskCreatedType
}
