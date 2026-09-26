package valueobject

import (
	"github.com/codememory1/d8r/internal/domain"
	"github.com/codememory1/d8r/pkg/ddd"
)

var _ ddd.ValueObject[ByteSize] = ByteSize{}

// ByteSize represents a validated resource size in bytes.
type ByteSize struct {
	value int64
}

// NewByteSize creates a byte size from a non-negative integer value.
func NewByteSize(value int64) (ByteSize, error) {
	if value < 0 {
		return ByteSize{}, domain.NewValidationError("invalid byte size", nil)
	}

	return ByteSize{value}, nil
}

// Int64 returns the size as a number of bytes.
func (v ByteSize) Int64() int64 {
	return v.value
}

// Equal reports whether two byte sizes are equal.
func (v ByteSize) Equal(other ByteSize) bool {
	return v.value == other.value
}
