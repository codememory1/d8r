package valueobject

import (
	"github.com/codememory1/d8r/internal/domain"
	"github.com/codememory1/d8r/pkg/ddd"
)

var _ ddd.ValueObject[DownloadStrategy] = DownloadStrategy{}

// DownloadStrategy identifies the strategy used to download a resource.
type DownloadStrategy struct {
	value string
}

const (
	downloadStrategySingle   = "single"
	downloadStrategyParallel = "parallel"
	downloadStrategyStream   = "stream"
)

// NewDownloadStrategy creates a validated download strategy from its string
// representation.
func NewDownloadStrategy(value string) (DownloadStrategy, error) {
	switch value {
	case downloadStrategySingle,
		downloadStrategyParallel,
		downloadStrategyStream:
		return DownloadStrategy{value}, nil
	default:
		return DownloadStrategy{}, domain.NewValidationError("invalid download strategy", nil)
	}
}

// SingleDownloadStrategy returns the sequential download strategy.
func SingleDownloadStrategy() DownloadStrategy {
	return DownloadStrategy{downloadStrategySingle}
}

// ParallelDownloadStrategy returns the parallel range download strategy.
func ParallelDownloadStrategy() DownloadStrategy {
	return DownloadStrategy{downloadStrategyParallel}
}

// StreamDownloadStrategy returns the streaming download strategy.
func StreamDownloadStrategy() DownloadStrategy {
	return DownloadStrategy{downloadStrategyStream}
}

// String returns the string representation of the download strategy.
func (v DownloadStrategy) String() string {
	return v.value
}

// Equal reports whether two download strategies are equal.
func (v DownloadStrategy) Equal(other DownloadStrategy) bool {
	return v.value == other.value
}
