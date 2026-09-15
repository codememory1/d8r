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

type Priority struct {
	value int
}

func NewPriority(priority int) (Priority, error) {
	if priority < minPriority || priority > maxPriority {
		return Priority{}, domain.NewValidationError("invalid priority range", nil)
	}

	return Priority{value: priority}, nil
}

func DefaultPriority() Priority {
	return Priority{value: defaultPriority}
}

func (v Priority) Int() int {
	return v.value
}

func (v Priority) Equal(other Priority) bool {
	return v.value == other.value
}
