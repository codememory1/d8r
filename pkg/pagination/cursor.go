package pagination

type Cursor struct {
	LastID    string `json:"last_id"`
	Timestamp int64  `json:"timestamp"`
}
