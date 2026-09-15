package ddd

// AggregateVersion tracks the version of an aggregate for optimistic concurrency control.
type AggregateVersion struct {
	value int64
}

// NewAggregateVersion restores an aggregate version from its persisted value.
func NewAggregateVersion(value int64) AggregateVersion {
	return AggregateVersion{
		value: value,
	}
}

// Version returns the current aggregate version.
func (v *AggregateVersion) Version() int64 {
	return v.value
}

// IncrementVersion advances the aggregate version by one.
func (v *AggregateVersion) IncrementVersion() {
	v.value++
}
