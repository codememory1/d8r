package entity

import (
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/ddd"
)

type TaskStatus string

const (
	TaskStatusPending     TaskStatus = "pending"
	TaskStatusDownloading TaskStatus = "downloading"
	TaskStatusCompleted   TaskStatus = "completed"
	TaskStatusFailed      TaskStatus = "failed"
)

var _ ddd.Entity[valueobject.ID] = (*Task)(nil)

type Task struct {
	ddd.AggregateVersion

	id        valueobject.ID
	url       valueobject.URL
	headers   valueobject.Headers
	filename  *valueobject.Filename
	priority  valueobject.Priority
	status    TaskStatus
	createdAt time.Time
	updatedAt *time.Time
}

func NewTask(
	url valueobject.URL,
	headers valueobject.Headers,
	filename *valueobject.Filename,
	priority valueobject.Priority,
) *Task {
	return &Task{
		id:        valueobject.NewID(),
		url:       url,
		headers:   headers,
		filename:  filename,
		priority:  priority,
		status:    TaskStatusPending,
		createdAt: time.Now(),
	}
}

func UnmarshalTask(
	id valueobject.ID,
	url valueobject.URL,
	headers valueobject.Headers,
	filename *valueobject.Filename,
	priority valueobject.Priority,
	status string,
	version int64,
	createdAt time.Time,
	updatedAt *time.Time,
) *Task {
	return &Task{
		AggregateVersion: ddd.NewAggregateVersion(version),

		id:        id,
		url:       url,
		headers:   headers,
		filename:  filename,
		priority:  priority,
		status:    TaskStatus(status),
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (t *Task) ID() valueobject.ID {
	return t.id
}

func (t *Task) URL() valueobject.URL {
	return t.url
}

func (t *Task) Headers() valueobject.Headers {
	return t.headers
}

func (t *Task) Filename() *valueobject.Filename {
	return t.filename
}

func (t *Task) Priority() valueobject.Priority {
	return t.priority
}

func (t *Task) Status() TaskStatus {
	return t.status
}

func (t *Task) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Task) UpdatedAt() *time.Time {
	return t.updatedAt
}

func (t *Task) Downloading() {
	t.status = TaskStatusDownloading
}

func (t *Task) Complete() {
	t.status = TaskStatusCompleted
}

func (t *Task) Fail() {
	t.status = TaskStatusFailed
}
