package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/codememory1/d8r/internal/infrastructure/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionPool struct {
	pool *pgxpool.Pool
}

// NewConnectionPool creates a new connection pool and performs a ping.
func NewConnectionPool(ctx context.Context, config config.Postgres) (*ConnectionPool, error) {
	poolConfig, err := pgxpool.ParseConfig(buildConnectionString(config))

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
	return p.pool.Query(ctx, sql, args...)
}

// QueryRow executes an SQL query and returns a single result row.
func (p *ConnectionPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

// Exec executes an SQL query that does not return a result set.
func (p *ConnectionPool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return p.pool.Exec(ctx, sql, args...)
}

// RunTransaction executes the provided function within a transaction.
// Upon successful execution, the transaction is committed;
// in case of an error, it is rolled back.
func (p *ConnectionPool) RunTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := p.pool.Begin(ctx)

	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}

		return err
	}

	return tx.Commit(ctx)
}

func (p *ConnectionPool) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}

func buildConnectionString(config config.Postgres) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
}
