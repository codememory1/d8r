package outbox

import "errors"

var (
	ErrConcurrentModification = errors.New("outbox message was concurrently modified")
)
