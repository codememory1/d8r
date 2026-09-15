package inspect

import "mime"

// parseContentType parses a Content-Type header and returns its normalized
// media type without parameters such as charset.
func parseContentType(raw string) (string, error) {
	mediaType, _, err := mime.ParseMediaType(raw)

	if err != nil {
		return "", err
	}

	return mediaType, nil
}
