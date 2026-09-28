package payload

type TaskInspectionStarted struct {
	TaskID       string `json:"task_id"`
	InspectionID string `json:"inspection_id"`
}
