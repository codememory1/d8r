package entity

import (
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/ddd"
)

// TaskStatus represents the current lifecycle state of a task.
type TaskStatus string

const (
	// TaskStatusPending indicates that the task is waiting to be inspected.
	TaskStatusPending TaskStatus = "pending"

	// TaskStatusInspecting indicates that the task is currently being inspected.
	TaskStatusInspecting TaskStatus = "inspecting"

	// TaskStatusReadyToDownload indicates that the task has been inspected and is ready to be downloaded.
	TaskStatusReadyToDownload TaskStatus = "ready_to_download"

	// TaskStatusDownloading indicates that the task is currently being downloaded.
	TaskStatusDownloading TaskStatus = "downloading"

	// TaskStatusCompleted indicates that the task has completed successfully.
	TaskStatusCompleted TaskStatus = "completed"

	// TaskStatusFailed indicates that the task has failed during processing.
	TaskStatusFailed TaskStatus = "failed"
)

var _ ddd.Entity[valueobject.ID] = (*Task)(nil)

// Task represents a downloadable resource and its current lifecycle state.
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

// NewTask creates a new task in the pending state.
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

// UnmarshalTask restores a task from its persisted state.
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

// ID returns the unique identifier of the task.
func (t *Task) ID() valueobject.ID {
	return t.id
}

// URL returns the URL of the resource to download.
func (t *Task) URL() valueobject.URL {
	return t.url
}

// Headers returns the HTTP headers associated with the task.
func (t *Task) Headers() valueobject.Headers {
	return t.headers
}

// Filename returns the requested filename, if one was provided.
func (t *Task) Filename() *valueobject.Filename {
	return t.filename
}

// Priority returns the processing priority of the task.
func (t *Task) Priority() valueobject.Priority {
	return t.priority
}

// Status returns the current lifecycle status of the task.
func (t *Task) Status() TaskStatus {
	return t.status
}

// CreatedAt returns the time when the task was created.
func (t *Task) CreatedAt() time.Time {
	return t.createdAt
}

// UpdatedAt returns the time when the task was last updated.
func (t *Task) UpdatedAt() *time.Time {
	return t.updatedAt
}

// Inspecting transitions the task to the inspecting state.
func (t *Task) Inspecting() {
	t.status = TaskStatusInspecting
}

// ReadyToDownload marks the task as successfully inspected and ready for download.
func (t *Task) ReadyToDownload() {
	t.status = TaskStatusReadyToDownload
}

// Downloading marks the task as currently being downloaded.
func (t *Task) Downloading() {
	t.status = TaskStatusDownloading
}

// Complete marks the task as successfully completed.
func (t *Task) Complete() {
	t.status = TaskStatusCompleted
}

// Fail marks the task as failed.
func (t *Task) Fail() {
	t.status = TaskStatusFailed
}
