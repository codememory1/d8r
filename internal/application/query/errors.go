package query

import "errors"

var (
	// ErrTaskNotFound is returned when the requested task does not exist.
	ErrTaskNotFound = errors.New("task not found")
)
