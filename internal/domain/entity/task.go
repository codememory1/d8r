package entity

import (
	"errors"
	"fmt"
	"time"

	domainevent "github.com/codememory1/d8r/internal/domain/event"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/ddd"
	"github.com/codememory1/d8r/pkg/statemachine"
)

// TaskTransition represents a named task state transition.
type TaskTransition string

// TaskStatus represents the current lifecycle state of a task.
type TaskStatus string

const (
	// TaskTransitionInspect starts task inspection.
	TaskTransitionInspect TaskTransition = "inspect"

	// TaskTransitionReady marks the task as ready for download.
	TaskTransitionReady TaskTransition = "ready"

	// TaskTransitionDownload starts task downloading.
	TaskTransitionDownload TaskTransition = "download"

	// TaskTransitionComplete completes the task.
	TaskTransitionComplete TaskTransition = "complete"

	// TaskTransitionFail marks the task as failed.
	TaskTransitionFail TaskTransition = "fail"
)

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

// taskStateMachine defines allowed task status transitions.
var taskStateMachine = statemachine.NewMachine(
	statemachine.T(
		TaskTransitionInspect,
		[]TaskStatus{
			TaskStatusPending,
		},
		TaskStatusInspecting,
	),
	statemachine.T(
		TaskTransitionReady,
		[]TaskStatus{
			TaskStatusInspecting,
		},
		TaskStatusReadyToDownload,
	),
	statemachine.T(
		TaskTransitionDownload,
		[]TaskStatus{
			TaskStatusReadyToDownload,
		},
		TaskStatusDownloading,
	),
	statemachine.T(
		TaskTransitionComplete,
		[]TaskStatus{
			TaskStatusDownloading,
		},
		TaskStatusCompleted,
	),
	statemachine.T(
		TaskTransitionFail,
		[]TaskStatus{
			TaskStatusInspecting,
			TaskStatusDownloading,
		},
		TaskStatusFailed,
	),
)

var (
	ErrInvalidTaskStatus = errors.New("invalid task status")
)

var _ ddd.Entity[valueobject.ID] = (*Task)(nil)

// Task represents a downloadable resource and its current lifecycle state.
type Task struct {
	ddd.AggregateRoot
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
	task := &Task{
		id:        valueobject.NewID(),
		url:       url,
		headers:   headers,
		filename:  filename,
		priority:  priority,
		status:    TaskStatusPending,
		createdAt: time.Now(),
	}

	task.Raise(domainevent.NewTaskCreated(task.id, task.createdAt))

	return task
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
) (*Task, error) {
	taskStatus, err := ParseTaskStatus(status)

	if err != nil {
		return nil, err
	}

	return &Task{
		AggregateVersion: ddd.NewAggregateVersion(version),
		id:               id,
		url:              url,
		headers:          headers,
		filename:         filename,
		priority:         priority,
		status:           taskStatus,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}, nil
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

// Inspect transitions the task to the inspecting state.
func (t *Task) Inspect() error {
	return t.transition(TaskTransitionInspect)
}

// Ready transitions the task to the ready-to-download state.
func (t *Task) Ready() error {
	return t.transition(TaskTransitionReady)
}

// Download transitions the task to the downloading state.
func (t *Task) Download() error {
	return t.transition(TaskTransitionDownload)
}

// Complete transitions the task to the completed state.
func (t *Task) Complete() error {
	return t.transition(TaskTransitionComplete)
}

// Fail transitions the task to the failed state.
func (t *Task) Fail() error {
	return t.transition(TaskTransitionFail)
}

// transition applies the named state transition to the task.
func (t *Task) transition(name TaskTransition) error {
	newStatus, err := taskStateMachine.Transition(name, t.status)

	if err != nil {
		return err
	}

	t.status = newStatus
	t.updatedAt = new(time.Now())

	return nil
}

func ParseTaskStatus(value string) (TaskStatus, error) {
	status := TaskStatus(value)

	switch status {
	case TaskStatusPending,
		TaskStatusInspecting,
		TaskStatusReadyToDownload,
		TaskStatusDownloading,
		TaskStatusCompleted,
		TaskStatusFailed:
		return status, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidTaskStatus, value)
	}
}
