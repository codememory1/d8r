package postgres

// Database combines PostgreSQL query execution, transaction management,
// and connection lifecycle operations.
type Database interface {
	Executor
	Transaction

	Close()
}
