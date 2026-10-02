package payload

// TaskDownloadProgress represents the payload of a task.download.progress webhook event.
type TaskDownloadProgress struct {
	TaskID          string `json:"task_id"`
	DownloadedBytes int64  `json:"downloaded_bytes"`
	TotalBytes      *int64 `json:"total_bytes"`
	ActiveRequests  int64  `json:"active_requests"`
}
