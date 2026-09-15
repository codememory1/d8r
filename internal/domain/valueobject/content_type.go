package valueobject

import (
	"mime"
	"strings"

	"github.com/codememory1/d8r/internal/domain"
	"github.com/codememory1/d8r/pkg/ddd"
)

var _ ddd.ValueObject[ContentType] = ContentType{}

var (
	ErrInvalidContentType = domain.NewValidationError("invalid content type", nil)
)

type ContentType struct {
	value string
}

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

func (v ContentType) String() string {
	return v.value
}

func (v ContentType) Equal(other ContentType) bool {
	return v.value == other.value
}
