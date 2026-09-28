package repository

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/codememory1/d8r/internal/domain/entity"
	"github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/internal/infrastructure/persistence/postgres"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

var _ repository.WebhookRepository = (*WebhookRepository)(nil)

type webhookModel struct {
	ID        string            `db:"id"`
	URL       string            `db:"url"`
	Headers   map[string]string `db:"headers"`
	Status    string            `db:"status"`
	CreatedAt time.Time         `db:"created_at"`
	UpdatedAt *time.Time        `db:"updated_at"`
}

type webhookEventModel struct {
	ID        string    `db:"id"`
	WebhookID string    `db:"webhook_id"`
	EventType string    `db:"event_type"`
	CreatedAt time.Time `db:"created_at"`
}

// WebhookRepository provides PostgreSQL persistence for webhook aggregates.
type WebhookRepository struct {
	connection *postgres.ConnectionPool
}

// NewWebhookRepository creates a PostgreSQL-backed webhook repository.
func NewWebhookRepository(connection *postgres.ConnectionPool) *WebhookRepository {
	return &WebhookRepository{
		connection: connection,
	}
}

// GetByID returns a webhook aggregate by its identifier.
func (r *WebhookRepository) GetByID(ctx context.Context, id valueobject.ID) (*entity.Webhook, error) {
	sql, args, err := sq.
		Select("*").
		From("webhooks").
		Where(sq.Eq{"id": id.String()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	var model webhookModel

	err = pgxscan.Get(ctx, r.connection, &model, sql, args...)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	webhooks, err := r.loadWebhookAggregates(ctx, []webhookModel{model})

	if err != nil {
		return nil, err
	}

	return webhooks[0], nil
}

// FindEnabledByEventType returns all enabled webhook aggregates subscribed
// to the specified event type.
func (r *WebhookRepository) FindEnabledByEventType(ctx context.Context, eventType valueobject.WebhookEventType) ([]*entity.Webhook, error) {
	// Select enabled webhooks that have a subscription to the requested event type.
	webhooksSQL, args, err := sq.
		Select("*").
		From("webhooks AS w").
		Where(sq.Eq{"w.status": entity.WebhookStatusEnabled}).
		Where(sq.Expr(`
			EXISTS (
				SELECT 1
				FROM webhook_events we
				WHERE we.webhook_id = w.id AND we.event_type = ?
			)
		`, eventType.String())).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	// Load the matching webhook persistence models.
	var webhookModels []webhookModel

	if err := pgxscan.Select(ctx, r.connection, &webhookModels, webhooksSQL, args...); err != nil {
		return nil, err
	}

	return r.loadWebhookAggregates(ctx, webhookModels)
}

// Save persists a webhook and all its event subscriptions atomically.
func (r *WebhookRepository) Save(ctx context.Context, webhook *entity.Webhook) error {
	webhookSql, webhookArgs, err := sq.
		Insert("webhooks").
		SetMap(map[string]any{
			"id":         webhook.ID().String(),
			"url":        webhook.URL().String(),
			"headers":    webhook.Headers().Map(),
			"status":     string(webhook.Status()),
			"created_at": webhook.CreatedAt(),
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	eventsBuilder := sq.
		Insert("webhook_events").
		Columns("id", "webhook_id", "event_type", "created_at").
		PlaceholderFormat(sq.Dollar)

	for _, webhookSub := range webhook.Subscriptions() {
		eventsBuilder = eventsBuilder.Values(
			valueobject.NewID().String(),
			webhook.ID().String(),
			webhookSub.EventType().String(),
			webhook.CreatedAt(),
		)
	}

	eventsSql, eventsArgs, err := eventsBuilder.ToSql()

	if err != nil {
		return err
	}

	return r.connection.Run(ctx, func(ctx context.Context) error {
		if _, err := r.connection.Exec(ctx, webhookSql, webhookArgs...); err != nil {
			return err
		}

		if _, err := r.connection.Exec(ctx, eventsSql, eventsArgs...); err != nil {
			return err
		}

		return nil
	})
}

// loadWebhookAggregates loads all subscriptions for the provided webhook
// models and reconstructs complete domain aggregates.
func (r *WebhookRepository) loadWebhookAggregates(ctx context.Context, webhookModels []webhookModel) ([]*entity.Webhook, error) {
	if len(webhookModels) == 0 {
		return []*entity.Webhook{}, nil
	}

	// Collect webhook identifiers for loading their subscriptions.
	webhookIDs := make([]string, len(webhookModels))

	for i, model := range webhookModels {
		webhookIDs[i] = model.ID
	}

	// Load every subscription belonging to the selected webhooks.
	eventsSQL, args, err := sq.
		Select(
			"id",
			"webhook_id",
			"event_type",
			"created_at",
		).
		From("webhook_events").
		Where(sq.Eq{
			"webhook_id": webhookIDs,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var eventModels []webhookEventModel

	if err := pgxscan.Select(ctx, r.connection, &eventModels, eventsSQL, args...); err != nil {
		return nil, err
	}

	// Group subscriptions by their owning webhook.
	eventModelsByWebhook := make(map[string][]webhookEventModel, len(webhookModels))

	for _, model := range eventModels {
		eventModelsByWebhook[model.WebhookID] = append(eventModelsByWebhook[model.WebhookID], model)
	}

	// Reconstruct complete webhook aggregates.
	webhooks := make([]*entity.Webhook, len(webhookModels))

	for i, model := range webhookModels {
		webhook, err := model.toDomainEntity(eventModelsByWebhook[model.ID])

		if err != nil {
			return nil, err
		}

		webhooks[i] = webhook
	}

	return webhooks, nil
}

// toDomainEntity reconstructs a webhook aggregate from its persistence model
// and subscription models.
func (m webhookModel) toDomainEntity(subscriptionModels []webhookEventModel) (*entity.Webhook, error) {
	id, err := valueobject.ParseID(m.ID)

	if err != nil {
		return nil, err
	}

	url, err := valueobject.NewURL(m.URL)

	if err != nil {
		return nil, err
	}

	headers, err := valueobject.NewHeaders(m.Headers)

	if err != nil {
		return nil, err
	}

	subscriptions := make(map[valueobject.WebhookEventType]entity.WebhookSubscription, len(subscriptionModels))

	for _, subscriptionModel := range subscriptionModels {
		subscriptionEntity, err := subscriptionModel.toDomainEntity()

		if err != nil {
			return nil, err
		}

		subscriptions[subscriptionEntity.EventType()] = subscriptionEntity
	}

	return entity.UnmarshalWebhook(id, url, headers, m.Status, subscriptions, m.CreatedAt, m.UpdatedAt)
}

// toDomainEntity reconstructs a webhook subscription from its persistence model.
func (m webhookEventModel) toDomainEntity() (entity.WebhookSubscription, error) {
	eventType, err := valueobject.NewWebhookEventType(m.EventType)

	if err != nil {
		return entity.WebhookSubscription{}, err
	}

	return entity.NewWebhookSubscription(eventType, m.CreatedAt), nil
}
