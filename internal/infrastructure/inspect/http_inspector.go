package inspect

import (
	"context"
	"net/http"

	"github.com/codememory1/d8r/internal/application/inspect"
)

// HttpInspector inspects HTTP resources and determines the appropriate
// httpdownload strategy from response metadata and range support.
type HttpInspector struct {
	client  *http.Client
	options Options
}

// NewHttpInspector creates an HTTP resource inspector with the provided
// client and httpdownload configuration.
func NewHttpInspector(client *http.Client, options Options) *HttpInspector {
	options.ReservedHeaders = append([]string(nil), options.ReservedHeaders...)

	return &HttpInspector{
		client:  client,
		options: options,
	}
}

// Inspect sends a HEAD request and, when necessary, a one-byte range probe,
// then builds a normalized inspection result.
func (i *HttpInspector) Inspect(ctx context.Context, url string, headers map[string]string) (inspect.Result, error) {
	headResp, err := i.sendHEAD(ctx, url, headers)

	if err != nil {
		return inspect.Result{}, err
	}

	defer headResp.Body.Close()

	var probeRangeResp *http.Response

	if i.shouldProbeRange(headResp) {
		probeRangeResp, err = i.sendProbeRange(ctx, url, headers)

		if err != nil {
			return inspect.Result{}, err
		}

		defer probeRangeResp.Body.Close()
	}

	return i.buildResult(headResp, probeRangeResp)
}
