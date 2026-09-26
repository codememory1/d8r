package eventcodec

import (
	"encoding/json"

	"github.com/codememory1/d8r/pkg/ddd"
)

// Decoder deserializes persisted event payloads into domain events.
type Decoder struct{}

// NewDecoder creates an event payload decoder.
func NewDecoder() *Decoder {
	return &Decoder{}
}

// Decode deserializes JSON into the provided event instance.
func (d *Decoder) Decode(payload []byte, target ddd.Event) error {
	return json.Unmarshal(payload, target)
}
