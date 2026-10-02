package event

import (
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/ddd"
)

var _ ddd.Event = (*TaskDownloadProgress)(nil)

// TaskDownloadProgressType identifies task download progress events.
const TaskDownloadProgressType ddd.EventType = "task.download.progress"

// TaskDownloadProgress represents an event containing a snapshot of a task's download progress.
type TaskDownloadProgress struct {
	TaskID          valueobject.ID
	DownloadedBytes int64
	TotalBytes      *int64
	ActiveRequests  int64
}

// NewTaskDownloadProgress creates a task-download-progress domain event.
func NewTaskDownloadProgress(taskID valueobject.ID, downloadedBytes int64, totalBytes *int64, activeRequests int64) TaskDownloadProgress {
	return TaskDownloadProgress{
		TaskID:          taskID,
		DownloadedBytes: downloadedBytes,
		TotalBytes:      totalBytes,
		ActiveRequests:  activeRequests,
	}
}

// Type returns the task download progress event type.
func (TaskDownloadProgress) Type() ddd.EventType {
	return TaskDownloadProgressType
}
