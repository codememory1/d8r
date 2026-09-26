package httpdownload

import (
	"errors"
	"fmt"
)

// HTTPRange represents an inclusive HTTP byte range.
type HTTPRange struct {
	Start int64
	End   int64
}

// NewHTTPRange creates an inclusive HTTP byte range from the provided offsets.
func NewHTTPRange(start int64, end int64) HTTPRange {
	return HTTPRange{start, end}
}

// ResolveHTTPRange calculates the HTTP range—specifically, the "start" and "end" for a specific part of a resource.
func ResolveHTTPRange(resourceSize int64, partNumber int, maxParts int) (HTTPRange, error) {
	if maxParts <= 0 {
		return HTTPRange{}, errors.New("max parts must be greater than zero")
	}

	if partNumber < 0 || partNumber >= maxParts {
		return HTTPRange{}, errors.New("invalid part number")
	}

	if resourceSize < int64(maxParts) {
		return HTTPRange{}, errors.New("resource size is too small for parts count")
	}

	// The size of a single resource part is calculated.
	partSize := resourceSize / int64(maxParts)

	// The start byte is calculated.
	start := partSize * int64(partNumber)

	// The final byte is calculated
	end := start + partSize - 1

	// For the last part, set the final byte
	// equal to the resource's last byte.
	if partNumber == maxParts-1 {
		end = resourceSize - 1
	}

	return NewHTTPRange(start, end), nil
}

// HeaderValue returns start and end in string format (bytes=<start>-<end>) for inclusion in an HTTP header.
func (r HTTPRange) HeaderValue() string {
	return fmt.Sprintf("bytes=%d-%d", r.Start, r.End)
}
