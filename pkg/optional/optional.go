package optional

// Map transforms an optional value and returns nil when the source value is nil.
func Map[T any, R any](value *T, mapper func(T) R) *R {
	if value == nil {
		return nil
	}

	return new(mapper(*value))
}
