package valueobject

import (
	"mime"
	"strings"

	"github.com/codememory1/d8r/internal/domain"
	"github.com/codememory1/d8r/pkg/ddd"
)

var _ ddd.ValueObject[ContentType] = ContentType{}

var (
	// ErrInvalidContentType is returned when a media type cannot be parsed.
	ErrInvalidContentType = domain.NewValidationError("invalid content type", nil)
)

// ContentType represents a validated and normalized media type.
type ContentType struct {
	value string
}

// NewContentType parses and normalizes a media type.
func NewContentType(value string) (ContentType, error) {
	raw := strings.TrimSpace(value)

	if raw == "" {
		return ContentType{}, ErrInvalidContentType
	}

	mediaType, params, err := mime.ParseMediaType(raw)

	if err != nil {
		return ContentType{}, ErrInvalidContentType
	}

	normalized := mime.FormatMediaType(mediaType, params)

	if normalized == "" {
		return ContentType{}, ErrInvalidContentType
	}

	return ContentType{normalized}, nil
}

// String returns the normalized media type.
func (v ContentType) String() string {
	return v.value
}

// Equal reports whether two content types are equal.
func (v ContentType) Equal(other ContentType) bool {
	return v.value == other.value
}
