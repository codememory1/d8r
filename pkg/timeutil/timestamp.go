package timeutil

import "time"

// UnixTimestamp converts an optional time value to a Unix timestamp.
func UnixTimestamp(value *time.Time) *int64 {
	if value == nil {
		return nil
	}

	return new(value.Unix())
}
