package eventregistry

import (
	"errors"

	"github.com/codememory1/d8r/pkg/ddd"
)

var ErrUnknownEventType = errors.New("unknown event type")

type Factory func() ddd.Event

type Registry struct {
	factories map[ddd.EventType]Factory
}

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
