package repository

import (
	"context"

	"github.com/codememory1/d8r/internal/domain/entity"
	"github.com/codememory1/d8r/internal/domain/valueobject"
)

// TaskRepository defines persistence operations for task entities.
type TaskRepository interface {
	// GetByID returns a task by its identifier.
	GetByID(ctx context.Context, id valueobject.ID) (*entity.Task, error)

	// Save persists a task entity.
	Save(ctx context.Context, task *entity.Task) error

	// Update updates the task entity.
	Update(ctx context.Context, task *entity.Task) error
}
