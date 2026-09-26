package download

import "github.com/codememory1/d8r/internal/domain/valueobject"

// StrategySelector selects a download strategy based on the resource size
// and its support for parallel range requests.
type StrategySelector struct {
	minParallelSize int64
}

// NewStrategySelector creates a strategy selector with the minimum resource
// size required for parallel downloading.
func NewStrategySelector(minParallelSize int64) *StrategySelector {
	return &StrategySelector{
		minParallelSize: minParallelSize,
	}
}

// Select returns the download strategy appropriate for the inspected resource.
func (s *StrategySelector) Select(
	size *valueobject.ByteSize,
	supportsParallel bool,
) valueobject.DownloadStrategy {
	if size == nil {
		return valueobject.StreamDownloadStrategy()
	}

	if size.Int64() < s.minParallelSize {
		return valueobject.SingleDownloadStrategy()
	}

	if !supportsParallel {
		return valueobject.SingleDownloadStrategy()
	}

	return valueobject.ParallelDownloadStrategy()
}
