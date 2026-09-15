package repository

import (
	"context"
	"errors"

	"github.com/codememory1/d8r/internal/domain/entity"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/pagination"
)

var (
	// ErrTaskNotFound is returned when a task with the requested identifier does not exist.
	ErrTaskNotFound = errors.New("task not found")
)

// TaskRepository defines persistence operations for task entities.
type TaskRepository interface {
	// GetAllPaginated returns tasks using cursor-based pagination.
	GetAllPaginated(ctx context.Context, cursor *pagination.Cursor, limit int) ([]*entity.Task, error)

	// GetById returns a task by its identifier.
	GetById(ctx context.Context, id valueobject.ID) (*entity.Task, error)

	// Save persists a task entity.
	Save(ctx context.Context, task *entity.Task) error

	// Update updates the task entity.
	Update(ctx context.Context, task *entity.Task) error
}
