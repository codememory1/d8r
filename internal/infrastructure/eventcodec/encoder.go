package eventcodec

import (
	"encoding/json"

	"github.com/codememory1/d8r/pkg/ddd"
)

type Encoder struct{}

func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encode serializes an event into JSON.
func (e *Encoder) Encode(event ddd.Event) ([]byte, error) {
	return json.Marshal(event)
}
