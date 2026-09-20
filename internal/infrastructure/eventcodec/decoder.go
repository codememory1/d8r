package eventcodec

import (
	"encoding/json"
	"errors"

	"github.com/codememory1/d8r/pkg/ddd"
)

var (
	// ErrUnknownEventType is returned when no factory is registered for an event type.
	ErrUnknownEventType = errors.New("unknown event type")
)

// Factory creates an empty event instance for decoding.
type Factory func() ddd.Event

// Decoder converts serialized event payloads into domain events.
type Decoder struct {
	factories map[ddd.EventType]Factory
}

// NewDecoder creates a new event decoder.
func NewDecoder() *Decoder {
	return &Decoder{
		factories: make(map[ddd.EventType]Factory),
	}
}

// Register associates an event type with its factory.
func (d *Decoder) Register(eventType ddd.EventType, factory Factory) {
	d.factories[eventType] = factory
}

// Decode deserializes a payload into the event registered for the specified type.
func (d *Decoder) Decode(eventType ddd.EventType, payload []byte) (*ddd.Event, error) {
	factory, ok := d.factories[eventType]

	if !ok {
		return nil, ErrUnknownEventType
	}

	event := factory()

	if err := json.Unmarshal(payload, event); err != nil {
		return nil, err
	}

	return &event, nil
}
