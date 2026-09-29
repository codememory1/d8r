package payload

// TaskDownloadFailed represents the payload of a task.download.failed webhook event.
type TaskDownloadFailed struct {
	TaskID string `json:"task_id"`
}
