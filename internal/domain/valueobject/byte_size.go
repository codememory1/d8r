package valueobject

import (
	"github.com/codememory1/d8r/internal/domain"
	"github.com/codememory1/d8r/pkg/ddd"
)

var _ ddd.ValueObject[ByteSize] = ByteSize{}

type ByteSize struct {
	value int64
}

func NewByteSize(value int64) (ByteSize, error) {
	if value < 0 {
		return ByteSize{}, domain.NewValidationError("invalid byte size", nil)
	}

	return ByteSize{value}, nil
}

func (v ByteSize) Int64() int64 {
	return v.value
}

func (v ByteSize) Equal(other ByteSize) bool {
	return v.value == other.value
}
