package repository

import (
	"context"

	"github.com/codememory1/d8r/internal/domain/entity"
)

type WebhookRepository interface {
	Save(ctx context.Context, webhook *entity.Webhook) error
}
