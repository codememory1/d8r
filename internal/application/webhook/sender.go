package webhook

import (
	"context"

	"github.com/codememory1/d8r/internal/domain/entity"
)

// Sender sends webhook delivery payloads to their target endpoints.
type Sender interface {
	// Send posts a delivery payload to the target webhook and returns the
	// received HTTP status code, or nil when no HTTP response was received.
	Send(ctx context.Context, webhook *entity.Webhook, delivery *entity.WebhookDelivery) (status *int, err error)
}
