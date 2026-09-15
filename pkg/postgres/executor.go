package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Executor defines PostgreSQL query and command execution operations.
type Executor interface {
	// Query executes a query that returns multiple rows.
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)

	// QueryRow executes a query expected to return at most one row.
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row

	// Exec executes a command that does not return rows.
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}
