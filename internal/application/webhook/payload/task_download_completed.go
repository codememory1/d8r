package payload

// TaskDownloadCompleted represents the payload of a task.download.completed webhook event.
type TaskDownloadCompleted struct {
	TaskID string `json:"task_id"`
}
