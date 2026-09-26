package valueobject

import (
	"regexp"
	"strings"

	"github.com/codememory1/d8r/internal/domain"
	"github.com/codememory1/d8r/pkg/ddd"
)

var _ ddd.ValueObject[Filename] = Filename{}

// filenameRegexp allows a simple filename with an optional extension.
// Paths and special characters are intentionally not allowed.
var filenameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_-]+(\.[a-z]+)?$`)

// Filename represents a validated filename without directory components.
type Filename struct {
	value string
}

// NewFilename creates and validates a filename.
func NewFilename(filename string) (Filename, error) {
	if filename == "" {
		return Filename{}, domain.NewValidationError("filename is empty", nil)
	}

	name := strings.TrimSpace(filename)

	if !filenameRegexp.MatchString(name) {
		return Filename{}, domain.NewValidationError("invalid filename", nil)
	}

	if len(name) > 255 {
		return Filename{}, domain.NewValidationError("filename is too long", nil)
	}

	return Filename{name}, nil
}

// String returns the filename as a string.
func (v Filename) String() string {
	return v.value
}

// Equal reports whether two filenames are equal.
func (v Filename) Equal(other Filename) bool {
	return v.value == other.value
}
