package query

import (
	"context"
	"time"

	"github.com/codememory1/d8r/pkg/cqrs"
	"github.com/codememory1/d8r/pkg/pagination"
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
	Items      []TaskListItem
	NextCursor *pagination.Cursor
}

// TaskListItem represents a task returned as part of the task list query.
type TaskListItem struct {
	ID        string
	URL       string
	Headers   map[string]string
	Filename  *string
	Priority  int
	Status    string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// ListTasksHandler handles ListTasks queries.
type ListTasksHandler struct {
	taskReader TaskReader
}

// NewListTasksHandler creates a new ListTasksHandler.
func NewListTasksHandler(taskReader TaskReader) *ListTasksHandler {
	return &ListTasksHandler{
		taskReader: taskReader,
	}
}

// Handle retrieves a page of tasks and returns a cursor for the next page when available.
func (h *ListTasksHandler) Handle(ctx context.Context, q ListTasks) (ListTasksResult, error) {
	// Fetch one additional task to determine whether another page exists.
	tasks, err := h.taskReader.GetAllPaginated(ctx, q.Cursor, q.Limit+1)

	if err != nil {
		return ListTasksResult{}, err
	}

	// Determine whether another page exists and remove the extra task.
	hasNext := len(tasks) > q.Limit

	if hasNext {
		tasks = tasks[:q.Limit]
	}

	// Map task read models into query result items.
	result := ListTasksResult{
		Items: make([]TaskListItem, len(tasks)),
	}

	for i, task := range tasks {
		result.Items[i] = TaskListItem{
			ID:        task.ID,
			URL:       task.URL,
			Headers:   task.Headers,
			Filename:  task.Filename,
			Priority:  task.Priority,
			Status:    task.Status,
			CreatedAt: task.CreatedAt,
			UpdatedAt: task.UpdatedAt,
		}
	}

	// Build the next cursor, if available, and include it in the result.
	if hasNext {
		lastTask := tasks[len(tasks)-1]

		result.NextCursor = &pagination.Cursor{
			LastID:    lastTask.ID,
			Timestamp: lastTask.CreatedAt.UnixMicro(),
		}
	}

	return result, nil
}
