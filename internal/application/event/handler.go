package event

import (
	"context"

	"github.com/codememory1/d8r/pkg/ddd"
)

// Handler handles a dispatched event.
type Handler interface {
	// Handle processes the provided event.
	Handle(ctx context.Context, event ddd.Event) error
}
