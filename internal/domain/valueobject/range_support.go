package valueobject

import "github.com/codememory1/d8r/pkg/ddd"

var _ ddd.ValueObject[RangeSupport] = RangeSupport(0)

// RangeSupport represents the detected level of HTTP byte-range support.
type RangeSupport uint8

const (
	// RangeSupportUnknown indicates that range support could not be determined.
	RangeSupportUnknown RangeSupport = iota

	// RangeSupportUnsupported indicates that the server ignored or rejected
	// byte-range requests.
	RangeSupportUnsupported

	// RangeSupportSupported indicates that the server supports byte-range requests.
	RangeSupportSupported
)

// Equal reports whether two range-support values are equal.
func (v RangeSupport) Equal(other RangeSupport) bool {
	return v == other
}
