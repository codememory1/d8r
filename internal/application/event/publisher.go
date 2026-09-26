package event

import (
	"context"

	"github.com/codememory1/d8r/pkg/ddd"
)

// Publisher publishes domain events for asynchronous processing.
type Publisher interface {
	// Publish publishes the provided domain event.
	Publish(ctx context.Context, event ddd.Event) error
}
