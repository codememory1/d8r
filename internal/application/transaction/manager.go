package transaction

import "context"

// Manager defines a boundary for executing operations within a transaction.
type Manager interface {
	// Run executes fn with a transaction-bound context, committing on success
	// and rolling back when fn returns an error.
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}
