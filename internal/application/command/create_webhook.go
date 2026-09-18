package command

import (
	"context"

	"github.com/codememory1/d8r/internal/domain/entity"
	"github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/cqrs"
)

// Compile-time check that CreateWebhookHandler implements the required command handler.
var _ cqrs.CommandHandler[CreateWebhook, valueobject.ID] = (*CreateWebhookHandler)(nil)

// CreateWebhook requests the creation of a webhook subscribed to the specified events.
type CreateWebhook struct {
	URL        valueobject.URL
	Headers    valueobject.Headers
	EventTypes []valueobject.WebhookEventType
}

// CreateWebhookHandler handles webhook creation commands.
type CreateWebhookHandler struct {
	webhookRepository repository.WebhookRepository
}

// NewCreateWebhookHandler creates a handler for the CreateWebhook command.
func NewCreateWebhookHandler(webhookRepository repository.WebhookRepository) *CreateWebhookHandler {
	return &CreateWebhookHandler{
		webhookRepository: webhookRepository,
	}
}

// Handle creates and persists a webhook and returns its identifier.
func (h *CreateWebhookHandler) Handle(ctx context.Context, cmd CreateWebhook) (valueobject.ID, error) {
	webhook, err := entity.NewWebhook(cmd.URL, cmd.Headers, cmd.EventTypes)

	if err != nil {
		return valueobject.ID{}, err
	}

	if saveErr := h.webhookRepository.Save(ctx, webhook); saveErr != nil {
		return valueobject.ID{}, saveErr
	}

	return webhook.ID(), nil
}
