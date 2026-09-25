package response

import (
	"github.com/codememory1/d8r/internal/application/query"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/timeutil"
)

type Task struct {
	ID        string            `json:"id"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Filename  *string           `json:"filename"`
	Priority  int               `json:"priority"`
	Status    string            `json:"status"`
	CreatedAt int64             `json:"created_at"`
	UpdatedAt *int64            `json:"updated_at"`
}

type CreateTask struct {
	ID string `json:"id"`
}

func NewCreateTask(id valueobject.ID) CreateTask {
	return CreateTask{
		ID: id.String(),
	}
}

func FromGetTaskResult(result query.GetTaskResult) Task {
	return Task{
		ID:        result.ID,
		URL:       result.URL,
		Headers:   result.Headers,
		Filename:  result.Filename,
		Priority:  result.Priority,
		Status:    result.Status,
		CreatedAt: result.CreatedAt.Unix(),
		UpdatedAt: timeutil.UnixTimestamp(result.UpdatedAt),
	}
}

func FromTaskListItems(items []query.TaskListItem) []Task {
	result := make([]Task, len(items))

	for i, item := range items {
		result[i] = Task{
			ID:        item.ID,
			URL:       item.URL,
			Headers:   item.Headers,
			Filename:  item.Filename,
			Priority:  item.Priority,
			Status:    item.Status,
			CreatedAt: item.CreatedAt.Unix(),
			UpdatedAt: timeutil.UnixTimestamp(item.UpdatedAt),
		}
	}

	return result
}
