package valueobject

import (
	"github.com/codememory1/d8r/internal/domain"
	"github.com/codememory1/d8r/pkg/ddd"
)

var _ ddd.ValueObject[Priority] = Priority{}

const (
	minPriority     = -256
	maxPriority     = 256
	defaultPriority = 0
)

// Priority represents the processing priority of a task.
type Priority struct {
	value int
}

// NewPriority creates a priority within the supported range.
func NewPriority(priority int) (Priority, error) {
	if priority < minPriority || priority > maxPriority {
		return Priority{}, domain.NewValidationError("invalid priority range", nil)
	}

	return Priority{value: priority}, nil
}

// DefaultPriority returns the default task priority.
func DefaultPriority() Priority {
	return Priority{value: defaultPriority}
}

// Int returns the numeric priority value.
func (v Priority) Int() int {
	return v.value
}

// Equal reports whether two priorities are equal.
func (v Priority) Equal(other Priority) bool {
	return v.value == other.value
}
