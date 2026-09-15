package valueobject

import (
	"net/url"

	"github.com/codememory1/d8r/internal/domain"
	"github.com/codememory1/d8r/pkg/ddd"
)

var _ ddd.ValueObject[URL] = URL{}

type URL struct {
	value string
}

func NewURL(value string) (URL, error) {
	if value == "" {
		return URL{}, domain.NewValidationError("URL is empty", nil)
	}

	parsedURL, err := url.Parse(value)

	if err != nil {
		return URL{}, domain.NewValidationError("URL parse error", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return URL{}, domain.NewValidationError("URL is invalid scheme", nil)
	}

	if parsedURL.Host == "" {
		return URL{}, domain.NewValidationError("URL host is empty", nil)
	}

	return URL{value}, nil
}

func (v URL) Equal(other URL) bool {
	return v.value == other.value
}

func (v URL) String() string {
	return v.value
}
