package valueobject

import (
	"github.com/codememory1/d8r/internal/domain"
	"github.com/codememory1/d8r/pkg/ddd"
	"github.com/google/uuid"
)

var _ ddd.ValueObject[ID] = ID{}

type ID struct {
	value uuid.UUID
}

func NewID() ID {
	return ID{uuid.New()}
}

func ParseID(value string) (ID, error) {
	parsedUUID, err := uuid.Parse(value)

	if err != nil {
		return ID{}, domain.NewValidationError("invalid id", err)
	}

	if parsedUUID == uuid.Nil {
		return ID{}, domain.NewValidationError("invalid id", nil)
	}

	return ID{parsedUUID}, nil
}

func (v ID) Equal(other ID) bool {
	return v.value == other.value
}

func (v ID) UUID() uuid.UUID {
	return v.value
}

func (v ID) String() string {
	return v.value.String()
}
