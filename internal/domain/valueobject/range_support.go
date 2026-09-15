package valueobject

import "github.com/codememory1/d8r/pkg/ddd"

var _ ddd.ValueObject[RangeSupport] = RangeSupport(0)

type RangeSupport uint8

const (
	RangeSupportUnknown RangeSupport = iota
	RangeSupportUnsupported
	RangeSupportSupported
)

func (v RangeSupport) Equal(other RangeSupport) bool {
	return v == other
}
