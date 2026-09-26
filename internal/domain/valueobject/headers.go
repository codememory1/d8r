package valueobject

import (
	"fmt"
	"maps"

	"github.com/codememory1/d8r/internal/domain"
	"github.com/codememory1/d8r/pkg/ddd"
	"golang.org/x/net/http/httpguts"
)

var _ ddd.ValueObject[Headers] = Headers{}

// Headers represents a validated collection of HTTP request headers.
type Headers struct {
	value map[string]string
}

// NewHeaders creates a validated and defensively copied header collection.
func NewHeaders(headers map[string]string) (Headers, error) {
	if headers == nil {
		headers = make(map[string]string)
	}

	clonedHeaders := maps.Clone(headers)

	for k, v := range clonedHeaders {
		// Header names must conform to the HTTP token syntax.
		if !httpguts.ValidHeaderFieldName(k) {
			return Headers{}, domain.NewValidationError(
				fmt.Sprintf("invalid header field name: %s", k),
				nil,
			)
		}

		// Reject control characters and other invalid bytes in header values.
		if !httpguts.ValidHeaderFieldValue(v) {
			return Headers{}, domain.NewValidationError(
				fmt.Sprintf("invalid header value in %s", k),
				nil,
			)
		}
	}

	return Headers{clonedHeaders}, nil
}

// Map returns a copy of the underlying HTTP headers.
func (v Headers) Map() map[string]string {
	return maps.Clone(v.value)
}

// Equal reports whether two header collections contain the same values.
func (v Headers) Equal(other Headers) bool {
	return maps.Equal(v.value, other.value)
}
