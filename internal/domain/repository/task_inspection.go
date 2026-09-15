package repository

import (
	"context"

	"github.com/codememory1/d8r/internal/domain/entity"
	"github.com/codememory1/d8r/internal/domain/valueobject"
)

// TaskInspectionRepository defines persistence operations for task entities.
type TaskInspectionRepository interface {
	// GetLastByTaskID returns the latest analysis for a task based on its ID.
	GetLastByTaskID(ctx context.Context, id valueobject.ID) (*entity.TaskInspection, error)

	// Save persists a task entity.
	Save(ctx context.Context, entity *entity.TaskInspection) error
}
