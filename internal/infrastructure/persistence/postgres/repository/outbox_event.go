package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/codememory1/d8r/internal/domain/entity"
	"github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/internal/infrastructure/persistence/postgres"
	"github.com/codememory1/d8r/pkg/ddd"
	"github.com/georgysavva/scany/v2/pgxscan"
)

// outboxEventClaimByStatusSQL atomically claims outbox events by transitioning
// them from the requested status to the processing status.
const outboxEventClaimByStatusSQL = `
	WITH filtered_outbox_events AS (
		SELECT
			id
		FROM outbox_events 
		WHERE status = $1
		ORDER BY id ASC
		LIMIT $2
		FOR UPDATE SKIP LOCKED
	)
	UPDATE outbox_events oe
	SET status = $3,
		version = version + 1,
		updated_at = NOW()
	FROM filtered_outbox_events
	WHERE oe.id = filtered_outbox_events.id
	RETURNING oe.*
`

// outboxEventModel represents the persistence model of an outbox event.
type outboxEventModel struct {
	ID        string     `db:"id"`
	EventType string     `db:"event_type"`
	Payload   []byte     `db:"payload"`
	Status    string     `db:"status"`
	Version   int64      `db:"version"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// OutboxEventRepository provides PostgreSQL persistence for outbox events.
type OutboxEventRepository struct {
	connection *postgres.ConnectionPool
}

// NewOutboxEventRepository creates a PostgreSQL-backed outbox event repository.
func NewOutboxEventRepository(pool *postgres.ConnectionPool) *OutboxEventRepository {
	return &OutboxEventRepository{
		connection: pool,
	}
}

// ClaimPending atomically claims pending outbox events for processing.
func (r *OutboxEventRepository) ClaimPending(ctx context.Context, limit int) ([]*entity.OutboxEvent, error) {
	var models []outboxEventModel

	err := pgxscan.Select(ctx, r.connection, &models, outboxEventClaimByStatusSQL, []any{
		entity.OutboxEventStatusPending,
		limit,
		entity.OutboxEventStatusProcessing,
	}...)

	if err != nil {
		return nil, err
	}

	entities := make([]*entity.OutboxEvent, len(models))

	for i, model := range models {
		e, err := model.toDomainEntity()

		if err != nil {
			return nil, err
		}

		entities[i] = e
	}

	return entities, nil
}

// Update updates the task entity in PostgreSQL using optimistic locking.
func (r *OutboxEventRepository) Update(ctx context.Context, outboxEvent *entity.OutboxEvent) error {
	sql, args, err := sq.
		Update("outbox_events").
		Set("event_type", outboxEvent.EventType()).
		Set("payload", outboxEvent.Payload()).
		Set("status", outboxEvent.Status()).
		Set("version", outboxEvent.Version()+1).
		Set("updated_at", time.Now()).
		Where(sq.Eq{
			"id":      outboxEvent.ID().String(),
			"version": outboxEvent.Version(),
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	tag, err := r.connection.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return repository.ErrConcurrentModification
	}

	outboxEvent.IncrementVersion()

	return nil
}

// Save persists a new outbox event.
func (r *OutboxEventRepository) Save(ctx context.Context, outboxEvent *entity.OutboxEvent) error {
	sql, args, err := sq.
		Insert("outbox_events").
		SetMap(map[string]any{
			"id":         outboxEvent.ID().String(),
			"event_type": outboxEvent.EventType(),
			"payload":    outboxEvent.Payload(),
			"status":     outboxEvent.Status(),
			"version":    outboxEvent.Version(),
			"created_at": outboxEvent.CreatedAt(),
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	if _, err := r.connection.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}

// toDomainEntity restores an outbox event from its persistence model.
func (m outboxEventModel) toDomainEntity() (*entity.OutboxEvent, error) {
	id, err := valueobject.ParseID(m.ID)

	if err != nil {
		return nil, err
	}

	return entity.UnmarshalOutboxEvent(
		id,
		ddd.EventType(m.EventType),
		m.Payload,
		m.Status,
		m.Version,
		m.CreatedAt,
		m.UpdatedAt,
	)
}
