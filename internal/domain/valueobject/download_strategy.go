package valueobject

import (
	"github.com/codememory1/d8r/internal/domain"
	"github.com/codememory1/d8r/pkg/ddd"
)

var _ ddd.ValueObject[DownloadStrategy] = DownloadStrategy{}

type DownloadStrategy struct {
	value string
}

const (
	downloadStrategySingle   = "single"
	downloadStrategyParallel = "parallel"
	downloadStrategyStream   = "stream"
)

func NewDownloadStrategy(value string) (DownloadStrategy, error) {
	switch value {
	case downloadStrategySingle,
		downloadStrategyParallel,
		downloadStrategyStream:
		return DownloadStrategy{value}, nil
	default:
		return DownloadStrategy{}, domain.NewValidationError("invalid httpdownload strategy", nil)
	}
}

func SingleDownloadStrategy() DownloadStrategy {
	return DownloadStrategy{downloadStrategySingle}
}

func ParallelDownloadStrategy() DownloadStrategy {
	return DownloadStrategy{downloadStrategyParallel}
}

func StreamDownloadStrategy() DownloadStrategy {
	return DownloadStrategy{downloadStrategyStream}
}

func (v DownloadStrategy) String() string {
	return v.value
}

func (v DownloadStrategy) Equal(other DownloadStrategy) bool {
	return v.value == other.value
}
