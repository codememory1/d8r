package inspect

import (
	"context"
	"net/http"
)

// sendHEAD sends a HEAD request to retrieve resource metadata without
// downloading the response body.
func (i *HttpInspector) sendHEAD(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	req, err := i.createRequest(ctx, url, http.MethodHead, headers)

	if err != nil {
		return nil, err
	}

	resp, err := i.client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// resolveHEADTotalSize returns the complete resource size advertised by a
// successful HEAD response. It returns nil when the size is unavailable.
func resolveHEADTotalSize(response *http.Response) *int64 {
	if response == nil {
		return nil
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil
	}

	if response.ContentLength < 0 {
		return nil
	}

	return &response.ContentLength
}
