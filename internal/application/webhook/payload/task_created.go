package payload

// TaskCreatedPayload represents the payload of a task.created webhook event.
type TaskCreatedPayload struct {
	TaskID string `json:"task_id"`
}
