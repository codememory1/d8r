package outbox

import "github.com/codememory1/d8r/pkg/ddd"

// Decoder defines a contract for decoding persisted payloads into domain
// events.
type Decoder interface {
	// Decode deserializes a payload into the provided event instance.
	Decode(payload []byte, target ddd.Event) error
}
