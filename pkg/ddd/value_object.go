package ddd

// ValueObject represents an immutable domain value compared by its contents.
type ValueObject[T any] interface {
	Equatable[T]
}
