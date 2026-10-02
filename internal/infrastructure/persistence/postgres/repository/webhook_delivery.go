package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/codememory1/d8r/internal/domain/entity"
	"github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/internal/infrastructure/persistence/postgres"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

var _ repository.WebhookDeliveryRepository = (*WebhookDeliveryRepository)(nil)

// webhookDeliveryModel represents a webhook delivery stored in PostgreSQL.
type webhookDeliveryModel struct {
	ID             string     `db:"id"`
	WebhookID      string     `db:"webhook_id"`
	SequenceKey    *string    `db:"sequence_key"`
	EventType      string     `db:"event_type"`
	Payload        []byte     `db:"payload"`
	Status         string     `db:"status"`
	Attempts       int        `db:"attempts"`
	NextAttemptAt  time.Time  `db:"next_attempt_at"`
	ResponseStatus *int       `db:"response_status"`
	LastError      *string    `db:"last_error"`
	Version        int64      `db:"version"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
	DeliveredAt    *time.Time `db:"delivered_at"`
}

// WebhookDeliveryRepository provides PostgreSQL persistence for webhook
// delivery aggregates.
type WebhookDeliveryRepository struct {
	connection *postgres.ConnectionPool
}

// NewWebhookDeliveryRepository creates a PostgreSQL-backed webhook delivery
// repository.
func NewWebhookDeliveryRepository(connection *postgres.ConnectionPool) *WebhookDeliveryRepository {
	return &WebhookDeliveryRepository{
		connection: connection,
	}
}

// GetByID returns a WebhookDelivery entity by ID.
func (r *WebhookDeliveryRepository) GetByID(ctx context.Context, id valueobject.ID) (*entity.WebhookDelivery, error) {
	sql, args, err := sq.
		Select("*").
		From("webhook_deliveries").
		Where(sq.Eq{"id": id.String()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var model webhookDeliveryModel

	err = pgxscan.Get(ctx, r.connection, &model, sql, args...)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return model.toDomainEntity()
}

// Save persists a new webhook delivery.
func (r *WebhookDeliveryRepository) Save(ctx context.Context, delivery *entity.WebhookDelivery) error {
	sql, args, err := sq.
		Insert("webhook_deliveries").
		SetMap(map[string]any{
			"id":              delivery.ID().String(),
			"webhook_id":      delivery.WebhookID().String(),
			"sequence_key":    delivery.SequenceKey(),
			"event_type":      delivery.EventType().String(),
			"payload":         delivery.Payload(),
			"status":          string(delivery.Status()),
			"attempts":        delivery.Attempts(),
			"next_attempt_at": delivery.NextAttemptAt(),
			"created_at":      delivery.CreatedAt(),
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("build webhook delivery insert: %w", err)
	}

	if _, err := r.connection.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("insert webhook delivery: %w", err)
	}

	return nil
}

// Update persists the current state of an existing webhook delivery.
func (r *WebhookDeliveryRepository) Update(ctx context.Context, delivery *entity.WebhookDelivery) error {
	sql, args, err := sq.
		Update("webhook_deliveries").
		SetMap(map[string]any{
			"status":          string(delivery.Status()),
			"attempts":        delivery.Attempts(),
			"next_attempt_at": delivery.NextAttemptAt(),
			"response_status": delivery.ResponseStatus(),
			"last_error":      delivery.LastError(),
			"version":         delivery.Version() + 1,
			"updated_at":      delivery.UpdatedAt(),
			"delivered_at":    delivery.DeliveredAt(),
		}).
		Where(sq.Eq{
			"id":      delivery.ID().String(),
			"version": delivery.Version(),
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("build webhook delivery update: %w", err)
	}

	result, err := r.connection.Exec(ctx, sql, args...)

	if err != nil {
		return fmt.Errorf("update webhook delivery: %w", err)
	}

	if result.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *WebhookDeliveryRepository) ClaimPending(ctx context.Context, limit int, maxAttempts int) ([]valueobject.ID, error) {
	var sql = `
		WITH filtered_webhook_deliveries AS (
			SELECT
				id
			FROM webhook_deliveries wd
			WHERE wd.status = @pending_status
				AND wd.next_attempt_at <= NOW()
				AND wd.attempts <= @attempts
				AND NOT EXISTS (
					SELECT
						1
					FROM webhook_deliveries wd2
					WHERE wd2.id < wd.id
						AND wd2.sequence_key = wd.sequence_key
						AND wd2.status IN (@pending_status, @processing_status)
				)
			ORDER BY id ASC
			LIMIT @limit
			FOR UPDATE SKIP LOCKED
		)
		UPDATE webhook_deliveries wd
		SET status = @processing_status,
			version = version + 1,
			updated_at = NOW()
		FROM filtered_webhook_deliveries fwd
		WHERE wd.id = fwd.id
		RETURNING wd.id
	`

	var models []webhookDeliveryModel

	err := pgxscan.Select(ctx, r.connection, &models, sql, pgx.NamedArgs{
		"pending_status":    entity.WebhookDeliveryStatusPending,
		"processing_status": entity.WebhookDeliveryStatusProcessing,
		"attempts":          maxAttempts,
		"limit":             limit,
	})

	if err != nil {
		return nil, err
	}

	ids := make([]valueobject.ID, len(models))

	for i, model := range models {
		id, err := valueobject.ParseID(model.ID)

		if err != nil {
			return nil, err
		}

		ids[i] = id
	}

	return ids, nil
}

// toDomainEntity reconstructs a webhook delivery domain entity from its
// persistence model.
func (m webhookDeliveryModel) toDomainEntity() (*entity.WebhookDelivery, error) {
	id, err := valueobject.ParseID(m.ID)

	if err != nil {
		return nil, err
	}

	webhookID, err := valueobject.ParseID(m.WebhookID)

	if err != nil {
		return nil, err
	}

	eventType, err := valueobject.NewWebhookEventType(m.EventType)

	if err != nil {
		return nil, err
	}

	return entity.UnmarshalWebhookDelivery(
		id,
		webhookID,
		m.SequenceKey,
		eventType,
		m.Payload,
		m.Status,
		m.Attempts,
		m.NextAttemptAt,
		m.ResponseStatus,
		m.LastError,
		m.Version,
		m.CreatedAt,
		m.UpdatedAt,
		m.DeliveredAt,
	)
}
