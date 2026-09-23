package outbox

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	infraoutbox "github.com/codememory1/d8r/internal/infrastructure/outbox"
	"github.com/codememory1/d8r/internal/infrastructure/persistence/postgres"
	"github.com/georgysavva/scany/v2/pgxscan"
)

var _ infraoutbox.Store = (*Store)(nil)

const claimPendingSQL = `
    WITH pending_messages AS (
        SELECT id
        FROM outbox_events
        WHERE status = $1
        ORDER BY created_at ASC, id ASC
        LIMIT $2
        FOR UPDATE SKIP LOCKED
    )
    UPDATE outbox_events oe
    SET
        status = $3,
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

type Store struct {
	connection *postgres.ConnectionPool
}

func NewStore(connection *postgres.ConnectionPool) *Store {
	return &Store{
		connection: connection,
	}
}

func (s *Store) ClaimPending(ctx context.Context, limit int) ([]infraoutbox.Message, error) {
	if limit <= 0 {
		return []infraoutbox.Message{}, nil
	}

	var rows []messageRow

	err := pgxscan.Select(
		ctx,
		s.connection,
		&rows,
		claimPendingSQL,
		infraoutbox.StatusPending,
		limit,
		infraoutbox.StatusProcessing,
	)

	if err != nil {
		return nil, err
	}

	messages := make([]infraoutbox.Message, len(rows))

	for i, row := range rows {
		messages[i] = row.toMessage()
	}

	return messages, nil
}

func (s *Store) MarkProcessed(ctx context.Context, message infraoutbox.Message) error {
	return s.mark(
		ctx,
		message,
		infraoutbox.StatusProcessed,
	)
}

func (s *Store) MarkFailed(ctx context.Context, message infraoutbox.Message) error {
	return s.mark(
		ctx,
		message,
		infraoutbox.StatusFailed,
	)
}

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
