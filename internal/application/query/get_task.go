package query

import (
	"context"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/cqrs"
	"github.com/codememory1/d8r/pkg/timeutil"
)

// Ensure GetTaskHandler implements the expected query handler contract.
var _ cqrs.QueryHandler[GetTask, GetTaskResult] = (*GetTaskHandler)(nil)

// GetTask identifies the task that should be retrieved.
type GetTask struct {
	ID valueobject.ID
}

// GetTaskResult contains the task data returned by GetTaskHandler.
type GetTaskResult struct {
	ID        string            `json:"id"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Filename  *string           `json:"filename"`
	Priority  int               `json:"priority"`
	Status    string            `json:"status"`
	CreatedAt int64             `json:"created_at"`
	UpdatedAt *int64            `json:"updated_at"`
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
		CreatedAt: task.CreatedAt.Unix(),
		UpdatedAt: timeutil.UnixTimestamp(task.UpdatedAt),
	}, nil
}
