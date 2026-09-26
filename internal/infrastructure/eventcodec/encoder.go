package eventcodec

import (
	"encoding/json"

	"github.com/codememory1/d8r/pkg/ddd"
)

// Encoder serializes domain events for persistence or transport.
type Encoder struct{}

// NewEncoder creates an event payload encoder.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encode serializes an event into JSON.
func (e *Encoder) Encode(event ddd.Event) ([]byte, error) {
	return json.Marshal(event)
}
