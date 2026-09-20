package ddd

// EventType identifies a domain event.
type EventType string

// Event represents an event that occurred in the domain.
type Event interface {
	Type() EventType
}
