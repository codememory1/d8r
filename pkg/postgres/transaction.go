package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Transaction defines transactional execution operations.
type Transaction interface {
	// RunTransaction executes fn within a transaction and commits it on success.
	// If fn returns an error, the transaction is rolled back.
	RunTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error
}
