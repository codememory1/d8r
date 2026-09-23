package outbox

import (
	"encoding/json"

	"github.com/codememory1/d8r/pkg/ddd"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusProcessed  Status = "processed"
	StatusFailed     Status = "failed"
)

type Message struct {
	ID        string
	EventType ddd.EventType
	Payload   json.RawMessage
	Version   int64
}
