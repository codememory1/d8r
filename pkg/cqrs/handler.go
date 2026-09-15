package cqrs

import "context"

// QueryHandler handles queries of type Q and returns results of type R.
type QueryHandler[Q, R any] interface {
	Handle(ctx context.Context, query Q) (R, error)
}

// CommandHandler handles commands of type C and returns results of type R.
type CommandHandler[C, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}
