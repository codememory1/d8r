package ddd

// Equatable defines equality comparison for values of type T.
type Equatable[T any] interface {
	Equal(other T) bool
}
