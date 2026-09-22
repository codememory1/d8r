package download

import "github.com/codememory1/d8r/internal/domain/valueobject"

type StrategySelector struct {
	minParallelSize int64
}

func NewStrategySelector(minParallelSize int64) *StrategySelector {
	return &StrategySelector{
		minParallelSize: minParallelSize,
	}
}

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
