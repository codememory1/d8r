package eventcodec

import (
	"encoding/json"

	"github.com/codememory1/d8r/pkg/ddd"
)

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

// Decode deserializes JSON into the provided event instance.
func (d *Decoder) Decode(payload []byte, target ddd.Event) error {
	return json.Unmarshal(payload, target)
}
