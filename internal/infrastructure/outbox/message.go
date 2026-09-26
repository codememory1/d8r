package outbox

import (
	"encoding/json"

	"github.com/codememory1/d8r/pkg/ddd"
)

// Status represents the processing state of an outbox message.
type Status string

const (
	// StatusPending indicates that the message is waiting to be processed.
	StatusPending Status = "pending"

	// StatusProcessing indicates that the message has been claimed for processing.
	StatusProcessing Status = "processing"

	// StatusProcessed indicates that the message was processed successfully.
	StatusProcessed Status = "processed"

	// StatusFailed indicates that message processing failed.
	StatusFailed Status = "failed"
)

// Message represents a persisted domain event claimed from the outbox.
type Message struct {
	ID        string
	EventType ddd.EventType
	Payload   json.RawMessage
	Version   int64
}
