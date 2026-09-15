package inspect

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidContentRange = errors.New("content range has invalid format")
	ErrInvalidRange        = errors.New("byte range has invalid format")

	ErrInvalidStart = errors.New("start offset is invalid")
	ErrInvalidEnd   = errors.New("end offset is invalid")
	ErrInvalidTotal = errors.New("total size is invalid")

	ErrNegativeStart = errors.New("start offset must be non-negative")
	ErrNegativeEnd   = errors.New("end offset must be non-negative")

	ErrInvalidRangeOrder = errors.New("start offset cannot be greater than end offset")
	ErrRangeExceedsTotal = errors.New("range exceeds total size")
)

// parseContentRange parses the Content-Range header value and returns the values of that header.
// The function expects a valid value for "bytes <start>-<end>/<total>".
func parseContentRange(value string) (int64, int64, *int64, error) {
	prefix := "bytes "

	// We verify that the value starts with the required prefix.
	// Example: "bytes 0-100/*" -> true
	// Example: "test 0-100/*" -> false
	if !strings.HasPrefix(value, prefix) {
		return 0, 0, nil, ErrInvalidContentRange
	}

	// Remove the prefix and extra spaces, normalizing the string.
	normalized := strings.TrimSpace(strings.TrimPrefix(value, prefix))

	// We divide the range and the total size into two parts.
	// Example: "0-100/300" -> "0-100", "300"
	rangePart, totalPart, ok := strings.Cut(normalized, "/")

	if !ok {
		return 0, 0, nil, ErrInvalidContentRange
	}

	// Split the range into start and end.
	// Example: "0-100" -> "0", "100"
	startPart, endPart, ok := strings.Cut(rangePart, "-")

	if !ok {
		return 0, 0, nil, ErrInvalidRange
	}

	// Parse the start of the range.
	// Example: "0" -> 0
	start, err := strconv.ParseInt(startPart, 10, 64)

	if err != nil {
		return 0, 0, nil, ErrInvalidStart
	}

	if start < 0 {
		return 0, 0, nil, ErrNegativeStart
	}

	// Parsing the end of the range
	// Example: "100" -> 100
	end, err := strconv.ParseInt(endPart, 10, 64)

	if err != nil {
		return 0, 0, nil, ErrInvalidEnd
	}

	if end < 0 {
		return 0, 0, nil, ErrNegativeEnd
	}

	if start > end {
		return 0, 0, nil, ErrInvalidRangeOrder
	}

	// According to RFC 9110, the total resource size can be "*";
	// therefore, if the resource size is unknown, it returns nil.
	if totalPart == "*" {
		return start, end, nil, nil
	}

	// Parsing the total size of the resource
	total, err := strconv.ParseInt(totalPart, 10, 64)

	if err != nil {
		return 0, 0, nil, ErrInvalidTotal
	}

	// If the resource size is specified, the size cannot be less than 0
	if total < 0 {
		return 0, 0, nil, ErrInvalidTotal
	}

	// The end of the range cannot extend beyond the limits of the resource.
	if end >= total {
		return 0, 0, nil, ErrRangeExceedsTotal
	}

	return start, end, &total, nil
}
