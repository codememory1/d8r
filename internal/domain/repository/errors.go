package repository

import "errors"

var (
	// ErrNotFound is returned when the requested entity does not exist.
	ErrNotFound = errors.New("not found")

	// ErrConcurrentModification is returned when an entity was changed
	// by another operation before the current update could be applied.
	ErrConcurrentModification = errors.New("entity was concurrently modified")
)
