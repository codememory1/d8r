package entity

import (
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
)

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

func (t *TaskInspection) ID() valueobject.ID {
	return t.id
}

func (t *TaskInspection) TaskID() valueobject.ID {
	return t.taskID
}

func (t *TaskInspection) EffectiveURL() valueobject.URL {
	return t.effectiveURL
}

func (t *TaskInspection) ContentType() *valueobject.ContentType {
	return t.contentType
}

func (t *TaskInspection) Filename() *valueobject.Filename {
	return t.filename
}

func (t *TaskInspection) Size() *valueobject.ByteSize {
	return t.size
}

func (t *TaskInspection) Strategy() valueobject.DownloadStrategy {
	return t.strategy
}

func (t *TaskInspection) ETag() *string {
	return t.etag
}

func (t *TaskInspection) LastModifiedAt() *time.Time {
	return t.lastModifiedAt
}

func (t *TaskInspection) CreatedAt() time.Time {
	return t.CreatedAt()
}
