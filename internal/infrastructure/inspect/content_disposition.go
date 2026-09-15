package inspect

import (
	"errors"
	"fmt"
	"mime"
	"strings"
)

var (
	// ErrNoContentDisposition is returned by the parseContentDisposition function
	// if the input value is an empty string.
	ErrNoContentDisposition = errors.New("no content disposition")

	// ErrInvalidContentDisposition is returned by the parseContentDisposition function
	// if the value fails to parse.
	ErrInvalidContentDisposition = errors.New("invalid content disposition")

	// ErrFilenameNotFound is returned by the parseContentDisposition function
	// if the filename could not be obtained after successfully parsing the value.
	ErrFilenameNotFound = errors.New("filename not found")
)

// parseContentDisposition parses the Content-Disposition header value and returns the value or an error.
func parseContentDisposition(raw string) (string, error) {
	// Normalize the string by removing extra spaces.
	normalized := strings.TrimSpace(raw)

	if normalized == "" {
		return "", ErrNoContentDisposition
	}

	_, params, err := mime.ParseMediaType(normalized)

	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidContentDisposition, err)
	}

	if filename := params["filename"]; filename != "" {
		return filename, nil
	}

	// Attempting to retrieve a filename containing Unicode characters defined in "filename*"
	// in accordance with RFC 5987 and RFC 8187.
	if filename := params["filename*"]; filename != "" {
		return filename, nil
	}

	// Failed to retrieve the file name
	return "", ErrFilenameNotFound
}
