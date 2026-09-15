package pagination

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

// base64Encoding is used to produce URL-safe opaque cursor values
// without padding characters.
var base64Encoding = base64.RawURLEncoding.Strict()

var (
	// ErrEmptyCursor is returned when an empty cursor value is provided.
	ErrEmptyCursor = errors.New("cursor is empty")
)

// Cursor represents the position of the last item from the previous page.
// It is used to fetch the next page using keyset pagination.
type Cursor struct {
	LastID    string `json:"last_id"`
	Timestamp int64  `json:"timestamp"`
}

// EncodeCursor serializes a cursor to JSON and encodes it
// as a URL-safe Base64 string suitable for use in query parameters.
func EncodeCursor(cursor Cursor) (string, error) {
	payload, err := json.Marshal(cursor)

	if err != nil {
		return "", err
	}

	return base64Encoding.EncodeToString(payload), nil
}

// DecodeCursor decodes a previously encoded cursor string
// back into its structured representation.
func DecodeCursor(raw string) (Cursor, error) {
	var cursor Cursor

	if raw == "" {
		return cursor, ErrEmptyCursor
	}

	payload, err := base64Encoding.DecodeString(raw)

	if err != nil {
		return cursor, err
	}

	if err = json.Unmarshal(payload, &cursor); err != nil {
		return cursor, err
	}

	return cursor, nil
}
