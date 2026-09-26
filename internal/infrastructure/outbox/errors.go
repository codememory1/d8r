package outbox

import "errors"

var (
	// ErrConcurrentModification is returned when an outbox message was changed
	// before its status update could be applied.
	ErrConcurrentModification = errors.New("outbox message was concurrently modified")
)
