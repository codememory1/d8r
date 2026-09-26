package eventregistry

import (
	"errors"

	"github.com/codememory1/d8r/pkg/ddd"
)

// ErrUnknownEventType is returned when no factory is registered for an event
// type.
var ErrUnknownEventType = errors.New("unknown event type")

// Factory creates an empty domain event instance for decoding.
type Factory func() ddd.Event

// Registry stores factories used to reconstruct domain events by type.
type Registry struct {
	factories map[ddd.EventType]Factory
}

// NewRegistry creates an empty domain event registry.
func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[ddd.EventType]Factory),
	}
}

// Register associates an event type with its factory.
func (r *Registry) Register(eventType ddd.EventType, factory Factory) {
	r.factories[eventType] = factory
}

// Get returns an empty event instance registered for the specified type.
func (r *Registry) Get(eventType ddd.EventType) (ddd.Event, error) {
	factory, ok := r.factories[eventType]

	if !ok {
		return nil, ErrUnknownEventType
	}

	return factory(), nil
}
