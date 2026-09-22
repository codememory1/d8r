package query

import (
	"context"
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/pagination"
)

// TaskReadModel represents a task optimized for read operations.
type TaskReadModel struct {
	ID        string
	URL       string
	Headers   map[string]string
	Filename  *string
	Priority  int
	Status    string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// TaskReader provides read-only access to tasks.
type TaskReader interface {
	// GetAllPaginated returns a paginated list of tasks.
	GetAllPaginated(ctx context.Context, cursor *pagination.Cursor, limit int) ([]TaskReadModel, error)

	// GetByID returns a task by its identifier.
	GetByID(ctx context.Context, id valueobject.ID) (TaskReadModel, error)
}
