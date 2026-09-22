package inspect

import (
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
)

// resolveHeader returns the first non-empty value of the requested header
// according to the provided response priority.
func resolveHeader(name string, responses ...*http.Response) string {
	for _, response := range responses {
		if response == nil {
			continue
		}

		value := strings.TrimSpace(response.Header.Get(name))

		if value != "" {
			return value
		}
	}

	return ""
}

// resolveEffectiveURL returns the final URL associated with the first
// available response, including redirects followed by the HTTP client.
func resolveEffectiveURL(responses ...*http.Response) string {
	for _, response := range responses {
		if response == nil || response.Request == nil || response.Request.URL == nil {
			continue
		}

		return response.Request.URL.String()
	}

	return ""
}

// resolveContentType returns the first valid Content-Type found in the
// provided responses.
func resolveContentType(responses ...*http.Response) *string {
	for _, response := range responses {
		if response == nil {
			continue
		}

		raw := response.Header.Get("Content-Type")

		if raw == "" {
			continue
		}

		contentType, err := parseContentType(raw)

		if err == nil {
			return &contentType
		}
	}

	return nil
}

// resolveFilename extracts a filename from Content-Disposition and falls back
// to the final URL path when the header does not provide one.
func resolveFilename(effectiveURL string, responses ...*http.Response) *string {
	for _, response := range responses {
		if response == nil {
			continue
		}

		raw := response.Header.Get("Content-Disposition")

		if raw == "" {
			continue
		}

		filename, err := parseContentDisposition(raw)

		if err == nil {
			return &filename
		}
	}

	parsedURL, err := url.Parse(effectiveURL)

	if err != nil {
		return nil
	}

	filename := path.Base(parsedURL.Path)

	if filename == "" || filename == "." || filename == "/" {
		return nil
	}

	filename, err = url.PathUnescape(filename)

	if err != nil {
		return nil
	}

	return &filename
}

// resolveTotalSize determines the complete resource size using Content-Range
// from a successful probe or Content-Length from a full response.
func resolveTotalSize(headResponse, probeRangeResponse *http.Response) *int64 {
	if probeRangeResponse == nil {
		return resolveHEADTotalSize(headResponse)
	}

	switch probeRangeResponse.StatusCode {
	case http.StatusOK:
		if probeRangeResponse.ContentLength > 0 {
			return &probeRangeResponse.ContentLength
		}

		return nil
	case http.StatusPartialContent:
		start, end, total, err := parseContentRange(probeRangeResponse.Header.Get("Content-Range"))

		if err != nil || start != 0 || end != 0 {
			return nil
		}

		return total
	default:
		return nil
	}
}

// resolveETag returns the first non-empty entity tag found in the provided
// responses.
func resolveETag(responses ...*http.Response) *string {
	value := resolveHeader("ETag", responses...)

	if value == "" {
		return nil
	}

	return &value
}

// resolveLastModified returns the first valid Last-Modified timestamp found
// in the provided responses.
func resolveLastModified(responses ...*http.Response) *time.Time {
	for _, response := range responses {
		if response == nil {
			continue
		}

		raw := response.Header.Get("Last-Modified")

		if raw == "" {
			continue
		}

		parsed, err := http.ParseTime(raw)

		if err == nil {
			return &parsed
		}
	}

	return nil
}

// resolveContentEncoding returns the content encoding of the representation
// that would be used for downloading. A nil value represents identity encoding.
func resolveContentEncoding(headResponse, probeRangeResponse *http.Response) *string {
	response := probeRangeResponse

	if response == nil {
		response = headResponse
	}

	if response == nil || response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil
	}

	encoding := strings.TrimSpace(response.Header.Get("Content-Encoding"))

	if encoding == "" {
		return nil
	}

	return &encoding
}

// resolveRangeSupport determines whether byte ranges were actually honored
// by validating the response to a one-byte range probe.
func resolveRangeSupport(probeRangeResponse *http.Response) valueobject.RangeSupport {
	if probeRangeResponse == nil {
		return valueobject.RangeSupportUnknown
	}

	switch probeRangeResponse.StatusCode {
	case http.StatusOK:
		return valueobject.RangeSupportUnsupported

	case http.StatusPartialContent:
		contentRange := probeRangeResponse.Header.Get("Content-Range")

		if contentRange == "" {
			return valueobject.RangeSupportUnknown
		}

		start, end, _, err := parseContentRange(contentRange)

		if err != nil {
			return valueobject.RangeSupportUnknown
		}

		if start != 0 || end != 0 {
			return valueobject.RangeSupportUnknown
		}

		return valueobject.RangeSupportSupported

	case http.StatusRequestedRangeNotSatisfiable:
		return valueobject.RangeSupportUnknown

	default:
		return valueobject.RangeSupportUnknown
	}
}

func resolveSupportsParallel(
	rangeSupport valueobject.RangeSupport,
	contentEncoding *string,
) bool {
	if rangeSupport == valueobject.RangeSupportSupported {
		return contentEncoding == nil || *contentEncoding == "" || strings.EqualFold(*contentEncoding, "identity")
	}

	return false
}
