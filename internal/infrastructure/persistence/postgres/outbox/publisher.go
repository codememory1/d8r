package outbox

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	appevent "github.com/codememory1/d8r/internal/application/event"
	infraoutbox "github.com/codememory1/d8r/internal/infrastructure/outbox"
	"github.com/codememory1/d8r/internal/infrastructure/persistence/postgres"
	"github.com/codememory1/d8r/pkg/ddd"
	"github.com/google/uuid"
)

var _ appevent.Publisher = (*Publisher)(nil)

type Publisher struct {
	connection *postgres.ConnectionPool
	encoder    infraoutbox.Encoder
}

func NewPublisher(
	connection *postgres.ConnectionPool,
	encoder infraoutbox.Encoder,
) *Publisher {
	return &Publisher{
		connection: connection,
		encoder:    encoder,
	}
}

func (p *Publisher) Publish(ctx context.Context, event ddd.Event) error {
	payload, err := p.encoder.Encode(event)

	if err != nil {
		return err
	}

	sql, args, err := sq.
		Insert("outbox_events").
		SetMap(map[string]any{
			"id":         uuid.NewString(),
			"event_type": event.Type(),
			"payload":    string(payload),
			"status":     infraoutbox.StatusPending,
			"created_at": time.Now(),
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("build outbox insert query: %w", err)
	}

	if _, err := p.connection.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("publish event to outbox: %w", err)
	}

	return nil
}
