package outbox

import "github.com/codememory1/d8r/pkg/ddd"

// EventRegistry provides event instances used to decode outbox payloads.
type EventRegistry interface {
	// Get returns a new event instance registered for the specified type.
	Get(eventType ddd.EventType) (ddd.Event, error)
}
