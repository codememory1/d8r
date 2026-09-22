package entity

import (
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
)

// TaskInspection represents the result of inspecting a task resource.
type TaskInspection struct {
	id             valueobject.ID
	taskID         valueobject.ID
	effectiveURL   valueobject.URL
	contentType    *valueobject.ContentType
	filename       *valueobject.Filename
	size           *valueobject.ByteSize
	strategy       valueobject.DownloadStrategy
	etag           *string
	lastModifiedAt *time.Time
	createdAt      time.Time
}

// NewTaskInspection creates a new task inspection.
func NewTaskInspection(
	taskId valueobject.ID,
	effectiveURL valueobject.URL,
	contentType *valueobject.ContentType,
	filename *valueobject.Filename,
	size *valueobject.ByteSize,
	strategy valueobject.DownloadStrategy,
	etag *string,
	lastModifiedAt *time.Time,
) *TaskInspection {
	return &TaskInspection{
		id:             valueobject.NewID(),
		taskID:         taskId,
		effectiveURL:   effectiveURL,
		contentType:    contentType,
		filename:       filename,
		size:           size,
		strategy:       strategy,
		etag:           etag,
		lastModifiedAt: lastModifiedAt,
		createdAt:      time.Now(),
	}
}

// UnmarshalTaskInspection restores a task inspection from persisted data.
func UnmarshalTaskInspection(
	id valueobject.ID,
	taskID valueobject.ID,
	effectiveURL valueobject.URL,
	contentType *valueobject.ContentType,
	filename *valueobject.Filename,
	size *valueobject.ByteSize,
	strategy valueobject.DownloadStrategy,
	etag *string,
	lastModifiedAt *time.Time,
	createdAt time.Time,
) *TaskInspection {
	return &TaskInspection{
		id:             id,
		taskID:         taskID,
		effectiveURL:   effectiveURL,
		contentType:    contentType,
		filename:       filename,
		size:           size,
		strategy:       strategy,
		etag:           etag,
		lastModifiedAt: lastModifiedAt,
		createdAt:      createdAt,
	}
}

// ID returns the inspection identifier.
func (t *TaskInspection) ID() valueobject.ID {
	return t.id
}

// TaskID returns the associated task identifier.
func (t *TaskInspection) TaskID() valueobject.ID {
	return t.taskID
}

// EffectiveURL returns the resolved resource URL.
func (t *TaskInspection) EffectiveURL() valueobject.URL {
	return t.effectiveURL
}

// ContentType returns the detected content type.
func (t *TaskInspection) ContentType() *valueobject.ContentType {
	return t.contentType
}

// Filename returns the detected filename.
func (t *TaskInspection) Filename() *valueobject.Filename {
	return t.filename
}

// Size returns the detected resource size.
func (t *TaskInspection) Size() *valueobject.ByteSize {
	return t.size
}

// Strategy returns the selected download strategy.
func (t *TaskInspection) Strategy() valueobject.DownloadStrategy {
	return t.strategy
}

// ETag returns the resource ETag.
func (t *TaskInspection) ETag() *string {
	return t.etag
}

// LastModifiedAt returns the resource's last modification time.
func (t *TaskInspection) LastModifiedAt() *time.Time {
	return t.lastModifiedAt
}

// CreatedAt returns the inspection creation time.
func (t *TaskInspection) CreatedAt() time.Time {
	return t.createdAt
}
