package inspect

import (
	"context"
	"net/http"
	"slices"
)

// createRequest creates an HTTP request and applies user-provided headers,
// excluding headers controlled internally by the inspector.
func (i *HttpInspector) createRequest(
	ctx context.Context,
	url string,
	method string,
	headers map[string]string,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)

	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		if slices.Contains(i.reservedHeaders, k) {
			continue
		}

		req.Header.Set(k, v)
	}

	return req, nil
}
