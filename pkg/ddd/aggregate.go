package ddd

// AggregateVersion tracks the version of an aggregate for optimistic concurrency control.
type AggregateVersion struct {
	value int64
}

type AggregateRoot struct {
	events []Event
}

// NewAggregateVersion restores an aggregate version from its persisted value.
func NewAggregateVersion(value int64) AggregateVersion {
	return AggregateVersion{
		value: value,
	}
}

func NewAggregateRoot() AggregateRoot {
	return AggregateRoot{}
}

// Version returns the current aggregate version.
func (v *AggregateVersion) Version() int64 {
	return v.value
}

// IncrementVersion advances the aggregate version by one.
func (v *AggregateVersion) IncrementVersion() {
	v.value++
}

func (v *AggregateRoot) Raise(event Event) {
	v.events = append(v.events, event)
}

func (v *AggregateRoot) PullEvents() []Event {
	events := make([]Event, len(v.events))

	copy(events, v.events)

	v.events = nil

	return events
}
