package event

import (
	"context"

	"github.com/codememory1/d8r/pkg/ddd"
)

type Publisher interface {
	Publish(ctx context.Context, event ddd.Event) error
}
