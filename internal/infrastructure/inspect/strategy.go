package inspect

import (
	"strings"

	"github.com/codememory1/d8r/internal/domain/valueobject"
)

// resolveStrategy selects stream for resources of unknown size, parallel for
// sufficiently large range-capable resources, and single in all other cases.
func (i *HttpInspector) resolveStrategy(
	totalSize *int64,
	rangeSupport valueobject.RangeSupport,
	contentEncoding *string,
) valueobject.DownloadStrategy {
	if totalSize == nil {
		return valueobject.StreamDownloadStrategy()
	}

	if *totalSize < i.config.MinParallelSize.Bytes() {
		return valueobject.SingleDownloadStrategy()
	}

	if rangeSupport != valueobject.RangeSupportSupported {
		return valueobject.SingleDownloadStrategy()
	}

	if contentEncoding != nil && *contentEncoding != "" && !strings.EqualFold(*contentEncoding, "identity") {
		return valueobject.SingleDownloadStrategy()
	}

	return valueobject.ParallelDownloadStrategy()
}
