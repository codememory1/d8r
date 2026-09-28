package valueobject

import (
	"encoding/json"

	"github.com/codememory1/d8r/internal/domain"
	"github.com/codememory1/d8r/pkg/ddd"
	"github.com/google/uuid"
)

var _ ddd.ValueObject[ID] = ID{}

// ID represents a non-empty UUID identifier.
type ID struct {
	value uuid.UUID
}

// NewID generates a new random identifier.
func NewID() ID {
	return ID{uuid.New()}
}

// ParseID parses and validates an identifier from its string representation.
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

// Equal reports whether two identifiers are equal.
func (v ID) Equal(other ID) bool {
	return v.value == other.value
}

// UUID returns the underlying UUID value.
func (v ID) UUID() uuid.UUID {
	return v.value
}

// String returns the canonical string representation of the identifier.
func (v ID) String() string {
	return v.value.String()
}

func (v ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

//goland:noinspection GoMixedReceiverTypes
func (v *ID) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	parsed, err := ParseID(value)

	if err != nil {
		return err
	}

	*v = parsed

	return nil
}
