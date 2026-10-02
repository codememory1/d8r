package outbox

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	infraoutbox "github.com/codememory1/d8r/internal/infrastructure/outbox"
	"github.com/codememory1/d8r/internal/infrastructure/persistence/postgres"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

var _ infraoutbox.Store = (*Store)(nil)

// Store provides PostgreSQL persistence operations for outbox messages.
type Store struct {
	connection *postgres.ConnectionPool
}

// NewStore creates a PostgreSQL-backed outbox store.
func NewStore(connection *postgres.ConnectionPool) *Store {
	return &Store{
		connection: connection,
	}
}

// ClaimPending atomically claims pending outbox messages using row locking.
func (s *Store) ClaimPending(ctx context.Context, limit int) ([]infraoutbox.Message, error) {
	if limit <= 0 {
		return []infraoutbox.Message{}, nil
	}

	var rows []messageRow

	sql := `
		WITH pending_messages AS (
			SELECT 
				oe.id
			FROM outbox_events oe
			WHERE status = @pending_status
				AND NOT EXISTS (
					SELECT
						1
					FROM outbox_events oe2
					WHERE oe2.id < oe.id
						AND oe2.sequence_key = oe.sequence_key
						AND oe2.status IN (@pending_status, @processing_status)
				)
			ORDER BY oe.id ASC
			LIMIT @limit
			FOR UPDATE SKIP LOCKED
		)
		UPDATE outbox_events oe
		SET
			status = @processing_status,
			version = version + 1,
			updated_at = NOW()
		FROM pending_messages
		WHERE oe.id = pending_messages.id
		RETURNING
			oe.id,
			oe.event_type,
			oe.payload,
			oe.version
	`

	err := pgxscan.Select(ctx, s.connection, &rows, sql, pgx.NamedArgs{
		"pending_status":    infraoutbox.StatusPending,
		"limit":             limit,
		"processing_status": infraoutbox.StatusProcessing,
	})

	if err != nil {
		return nil, err
	}

	messages := make([]infraoutbox.Message, len(rows))

	for i, row := range rows {
		messages[i] = row.toMessage()
	}

	return messages, nil
}

// MarkProcessed marks an outbox message as successfully processed.
func (s *Store) MarkProcessed(ctx context.Context, message infraoutbox.Message) error {
	return s.mark(
		ctx,
		message,
		infraoutbox.StatusProcessed,
	)
}

// MarkFailed marks an outbox message as failed.
func (s *Store) MarkFailed(ctx context.Context, message infraoutbox.Message) error {
	return s.mark(
		ctx,
		message,
		infraoutbox.StatusFailed,
	)
}

// mark updates an outbox message status using optimistic locking.
func (s *Store) mark(ctx context.Context, message infraoutbox.Message, status infraoutbox.Status) error {
	sql, args, err := sq.
		Update("outbox_events").
		Set("status", status).
		Set("version", sq.Expr("version + 1")).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{
			"id":      message.ID,
			"version": message.Version,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	tag, err := s.connection.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return infraoutbox.ErrConcurrentModification
	}

	return nil
}
