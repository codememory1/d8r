package eventbus

import (
	"context"

	appevent "github.com/codememory1/d8r/internal/application/event"
	"github.com/codememory1/d8r/pkg/ddd"
)

type Dispatcher struct {
	handlers map[ddd.EventType][]appevent.Handler
}

// NewDispatcher creates a dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[ddd.EventType][]appevent.Handler),
	}
}

// Subscribe registers a handler for the specified event type.
func (d *Dispatcher) Subscribe(eventType ddd.EventType, handler appevent.Handler) {
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// Dispatch passes the event to all handlers registered for its type.
func (d *Dispatcher) Dispatch(ctx context.Context, event ddd.Event) error {
	for _, handler := range d.handlers[event.Type()] {
		if err := handler.Handle(ctx, event); err != nil {
			return err
		}
	}

	return nil
}
