package outbox

import "github.com/codememory1/d8r/pkg/ddd"

// Encoder defines a contract for serializing domain events.
type Encoder interface {
	// Encode serializes the provided domain event.
	Encode(event ddd.Event) ([]byte, error)
}
