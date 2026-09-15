package query

import (
	"context"

	"github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/cqrs"
	"github.com/codememory1/d8r/pkg/optional"
	"github.com/codememory1/d8r/pkg/pagination"
	"github.com/codememory1/d8r/pkg/timeutil"
)

var _ cqrs.QueryHandler[ListTasks, ListTasksResult] = (*ListTasksHandler)(nil)

const (
	// ListTasksMinLimit defines the minimum allowed page size.
	ListTasksMinLimit = 1

	// ListTasksMaxLimit defines the maximum allowed page size.
	ListTasksMaxLimit = 100
)

// ListTasks represents a query for retrieving tasks using cursor-based pagination.
type ListTasks struct {
	Cursor *pagination.Cursor
	Limit  int
}

// ListTasksResult represents a paginated list of tasks.
type ListTasksResult struct {
	Items      []TaskListItem `json:"items"`
	NextCursor *string        `json:"next_cursor"`
}

// TaskListItem represents a task returned as part of the task list query.
type TaskListItem struct {
	ID        string            `json:"id"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Filename  *string           `json:"filename"`
	Priority  int               `json:"priority"`
	Status    string            `json:"status"`
	CreatedAt int64             `json:"created_at"`
	UpdatedAt *int64            `json:"updated_at"`
}

// ListTasksHandler handles ListTasks queries.
type ListTasksHandler struct {
	taskRepository repository.TaskRepository
}

// NewListTasksHandler creates a new ListTasksHandler.
func NewListTasksHandler(taskRepository repository.TaskRepository) *ListTasksHandler {
	return &ListTasksHandler{
		taskRepository: taskRepository,
	}
}

// Handle retrieves a page of tasks and returns a cursor for the next page when available.
func (h *ListTasksHandler) Handle(ctx context.Context, query ListTasks) (ListTasksResult, error) {
	tasks, err := h.taskRepository.GetAllPaginated(ctx, query.Cursor, query.Limit+1)

	if err != nil {
		return ListTasksResult{}, err
	}

	hasNext := len(tasks) > query.Limit

	if hasNext {
		tasks = tasks[:query.Limit]
	}

	result := ListTasksResult{
		Items: make([]TaskListItem, len(tasks)),
	}

	for i, task := range tasks {
		result.Items[i] = TaskListItem{
			ID:        task.ID().String(),
			URL:       task.URL().String(),
			Headers:   task.Headers().Map(),
			Filename:  optional.Map(task.Filename(), valueobject.Filename.String),
			Priority:  task.Priority().Int(),
			Status:    string(task.Status()),
			CreatedAt: task.CreatedAt().Unix(),
			UpdatedAt: timeutil.UnixTimestamp(task.UpdatedAt()),
		}
	}

	// Build the next cursor, if available, and include it in the result.
	if hasNext {
		lastTask := tasks[len(tasks)-1]
		nextCursor, err := pagination.EncodeCursor(pagination.Cursor{
			LastID:    lastTask.ID().String(),
			Timestamp: lastTask.CreatedAt().UnixMicro(),
		})

		if err != nil {
			return ListTasksResult{}, err
		}

		result.NextCursor = &nextCursor
	}

	return result, nil
}
