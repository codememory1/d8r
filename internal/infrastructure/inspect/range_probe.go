package inspect

import (
	"context"
	"net/http"
)

// sendProbeRange requests the first byte of a resource to verify whether
// the server honors byte-range requests.
func (i *HttpInspector) sendProbeRange(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	req, err := i.createRequest(ctx, url, http.MethodGet, headers)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Range", "bytes=0-0")

	resp, err := i.client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// shouldProbeRange reports whether a range probe is required because HEAD is
// unsupported, the resource size is unknown, or parallel downloading is possible.
func (i *HttpInspector) shouldProbeRange(headResp *http.Response) bool {
	if headResp == nil {
		return false
	}

	if headResp.StatusCode == http.StatusMethodNotAllowed ||
		headResp.StatusCode == http.StatusNotImplemented {
		return true
	}

	if headResp.StatusCode < 200 || headResp.StatusCode >= 300 {
		return false
	}

	if headResp.ContentLength < 0 {
		return true
	}

	return headResp.ContentLength >= i.config.MinParallelSize.Bytes()
}
