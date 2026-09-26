package query

import (
	"context"
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/cqrs"
)

// Ensure GetTaskHandler implements the expected query handler contract.
var _ cqrs.QueryHandler[GetTask, GetTaskResult] = (*GetTaskHandler)(nil)

// GetTask identifies the task that should be retrieved.
type GetTask struct {
	ID valueobject.ID
}

// GetTaskResult contains the task data returned by GetTaskHandler.
type GetTaskResult struct {
	ID        string
	URL       string
	Headers   map[string]string
	Filename  *string
	Priority  int
	Status    string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// GetTaskHandler handles queries for retrieving a task by its identifier.
type GetTaskHandler struct {
	taskReader TaskReader
}

// NewGetTaskHandler creates a handler for task retrieval queries.
func NewGetTaskHandler(taskReader TaskReader) *GetTaskHandler {
	return &GetTaskHandler{
		taskReader: taskReader,
	}
}

// Handle retrieves a task and converts it into a query result.
func (h *GetTaskHandler) Handle(ctx context.Context, q GetTask) (GetTaskResult, error) {
	task, err := h.taskReader.GetByID(ctx, q.ID)

	if err != nil {
		return GetTaskResult{}, err
	}

	return GetTaskResult{
		ID:        task.ID,
		URL:       task.URL,
		Headers:   task.Headers,
		Filename:  task.Filename,
		Priority:  task.Priority,
		Status:    task.Status,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}, nil
}
