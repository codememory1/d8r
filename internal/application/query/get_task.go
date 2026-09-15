package query

import (
	"context"

	"github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/cqrs"
	"github.com/codememory1/d8r/pkg/optional"
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
	taskRepository repository.TaskRepository
}

// NewGetTaskHandler creates a handler for task retrieval queries.
func NewGetTaskHandler(taskRepository repository.TaskRepository) *GetTaskHandler {
	return &GetTaskHandler{
		taskRepository: taskRepository,
	}
}

// Handle retrieves a task and converts it into a query result.
func (h *GetTaskHandler) Handle(ctx context.Context, query GetTask) (GetTaskResult, error) {
	task, err := h.taskRepository.GetById(ctx, query.ID)

	if err != nil {
		return GetTaskResult{}, err
	}

	return GetTaskResult{
		ID:        task.ID().String(),
		URL:       task.URL().String(),
		Headers:   task.Headers().Map(),
		Filename:  optional.Map(task.Filename(), valueobject.Filename.String),
		Status:    string(task.Status()),
		CreatedAt: task.CreatedAt().Unix(),
		UpdatedAt: timeutil.UnixTimestamp(task.UpdatedAt()),
	}, nil
}
