package inspect

import (
	"fmt"
	"net/http"

	"github.com/codememory1/d8r/internal/application/inspect"
	"github.com/codememory1/d8r/internal/domain/valueobject"
)

// buildResult converts inspected HTTP responses into a validated application
// result and determines the appropriate httpdownload strategy.
func (i *HttpInspector) buildResult(headResponse, probeRangeResponse *http.Response) (inspect.Result, error) {
	rawEffectiveURL := resolveEffectiveURL(probeRangeResponse, headResponse)
	effectiveURL, err := valueobject.NewURL(rawEffectiveURL)

	if err != nil {
		return inspect.Result{}, fmt.Errorf("create effective URL: %w", err)
	}

	var contentType *valueobject.ContentType

	rawContentType := resolveContentType(probeRangeResponse, headResponse)

	if rawContentType != nil {
		value, err := valueobject.NewContentType(*rawContentType)

		if err != nil {
			return inspect.Result{}, fmt.Errorf("create content type: %w", err)
		}

		contentType = &value
	}

	var filename *valueobject.Filename

	rawFilename := resolveFilename(
		rawEffectiveURL,
		probeRangeResponse,
		headResponse,
	)

	if rawFilename != nil {
		value, err := valueobject.NewFilename(*rawFilename)

		if err != nil {
			return inspect.Result{}, fmt.Errorf("create filename: %w", err)
		}

		filename = &value
	}

	totalSize := resolveTotalSize(headResponse, probeRangeResponse)

	var byteSize *valueobject.ByteSize

	if totalSize != nil {
		value, err := valueobject.NewByteSize(*totalSize)

		if err != nil {
			return inspect.Result{}, fmt.Errorf("create byte size: %w", err)
		}

		byteSize = &value
	}

	rangeSupport := resolveRangeSupport(probeRangeResponse)
	contentEncoding := resolveContentEncoding(headResponse, probeRangeResponse)

	return inspect.Result{
		EffectiveURL:     effectiveURL,
		ContentType:      contentType,
		Filename:         filename,
		Size:             byteSize,
		DownloadStrategy: i.resolveStrategy(totalSize, rangeSupport, contentEncoding),
		ETag:             resolveETag(probeRangeResponse, headResponse),
		LastModified:     resolveLastModified(probeRangeResponse, headResponse),
	}, nil
}
