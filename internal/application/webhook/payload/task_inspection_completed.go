package payload

// TaskInspectionCompleted is the webhook payload sent after a task
// inspection completes successfully.
type TaskInspectionCompleted struct {
	TaskID       string `json:"task_id"`
	InspectionID string `json:"inspection_id"`
}
