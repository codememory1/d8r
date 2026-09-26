package postgres

import (
	"context"
	"fmt"

	"github.com/codememory1/d8r/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

const (
	transactionContextKey contextKey = "transaction"
)

type ConnectionPool struct {
	pool *pgxpool.Pool
}

// NewConnectionPool creates a new connection pool and performs a ping.
func NewConnectionPool(ctx context.Context, options Options) (*ConnectionPool, error) {
	poolConfig, err := pgxpool.ParseConfig(buildConnectionString(options))

	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)

	if err != nil {
		return nil, err
	}

	// Pinging the database to test the connection.
	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		return nil, err
	}

	return &ConnectionPool{pool}, nil
}

// Query executes an SQL query and returns a result set of rows.
func (p *ConnectionPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.executor(ctx).Query(ctx, sql, args...)
}

// QueryRow executes an SQL query and returns a single result row.
func (p *ConnectionPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.executor(ctx).QueryRow(ctx, sql, args...)
}

// Exec executes an SQL query that does not return a result set.
func (p *ConnectionPool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return p.executor(ctx).Exec(ctx, sql, args...)
}

// Run executes fn within a database transaction and provides a context bound
// to that transaction. It commits on success and rolls back on error.
func (p *ConnectionPool) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	return pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		txCtx := context.WithValue(ctx, transactionContextKey, tx)

		return fn(txCtx)
	})
}

// Close releases all database connections managed by the pool.
func (p *ConnectionPool) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}

// executor returns the transaction associated with ctx or the underlying pool.
func (p *ConnectionPool) executor(ctx context.Context) postgres.Executor {
	if tx, ok := ctx.Value(transactionContextKey).(pgx.Tx); ok {
		return tx
	}

	return p.pool
}

// buildConnectionString constructs a PostgreSQL connection URL from the configuration.
func buildConnectionString(options Options) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		options.User,
		options.Password,
		options.Host,
		options.Port,
		options.Database,
	)
}
