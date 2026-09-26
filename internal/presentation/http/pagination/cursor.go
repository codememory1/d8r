package pagination

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	corepagination "github.com/codememory1/d8r/pkg/pagination"
)

var (
	// ErrEmptyCursor is returned when an empty pagination cursor is provided.
	ErrEmptyCursor = errors.New("cursor is empty")
	base64Encoding = base64.RawURLEncoding.Strict()
)

type cursorPayload struct {
	LastID    string `json:"last_id"`
	Timestamp int64  `json:"timestamp"`
}

// Encode serializes a pagination cursor into an opaque URL-safe string.
func Encode(cursor corepagination.Cursor) (string, error) {
	payload, err := json.Marshal(cursorPayload{
		LastID:    cursor.LastID,
		Timestamp: cursor.Timestamp,
	})

	if err != nil {
		return "", err
	}

	return base64Encoding.EncodeToString(payload), nil
}

// Decode restores a pagination cursor from an opaque URL-safe string.
func Decode(raw string) (corepagination.Cursor, error) {
	if raw == "" {
		return corepagination.Cursor{}, ErrEmptyCursor
	}

	payload, err := base64Encoding.DecodeString(raw)
	if err != nil {
		return corepagination.Cursor{}, err
	}

	var decoded cursorPayload

	if err := json.Unmarshal(payload, &decoded); err != nil {
		return corepagination.Cursor{}, err
	}

	return corepagination.Cursor{
		LastID:    decoded.LastID,
		Timestamp: decoded.Timestamp,
	}, nil
}
