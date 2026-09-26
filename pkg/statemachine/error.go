package statemachine

import (
	"errors"
	"fmt"
)

var (
	// ErrTransitionNotFound indicates that the requested transition is not registered.
	ErrTransitionNotFound = errors.New("transition not found")
)

// InvalidTransitionError indicates that a transition cannot be applied
// from the current state.
type InvalidTransitionError[N comparable, S comparable] struct {
	Name N
	From S
}

// Error implements the error interface.
func (e InvalidTransitionError[N, S]) Error() string {
	return fmt.Sprintf(
		"transition %v is not allowed from state %v",
		e.Name,
		e.From,
	)
}
