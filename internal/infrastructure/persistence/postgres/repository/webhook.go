package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/codememory1/d8r/internal/domain/entity"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/internal/infrastructure/persistence/postgres"
)

type WebhookRepository struct {
	connection *postgres.ConnectionPool
}

func NewWebhookRepository(connection *postgres.ConnectionPool) *WebhookRepository {
	return &WebhookRepository{
		connection: connection,
	}
}

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
			webhookSub.EventType(),
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
