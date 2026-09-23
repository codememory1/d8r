package event

import (
	"context"

	"github.com/codememory1/d8r/pkg/ddd"
)

// Dispatcher manages event subscriptions and dispatches events to registered handlers.
type Dispatcher interface {
	// Dispatch invokes all handlers registered for the specified event type.
	Dispatch(ctx context.Context, event ddd.Event) error
}
