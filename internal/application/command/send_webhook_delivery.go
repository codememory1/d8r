package command

import (
	"context"
	"errors"
	"time"

	"github.com/codememory1/d8r/internal/application/webhook"
	"github.com/codememory1/d8r/internal/domain/entity"
	"github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/cqrs"
)

var _ cqrs.CommandHandler[SendWebhookDelivery, struct{}] = (*SendWebhookDeliveryHandler)(nil)

type SendWebhookDelivery struct {
	WebhookDeliveryID valueobject.ID
}

type SendWebhookDeliveryHandler struct {
	webhookRepository         repository.WebhookRepository
	webhookDeliveryRepository repository.WebhookDeliveryRepository
	sender                    webhook.Sender
	maxAttempts               int
	retryDelay                time.Duration
}

func NewSendWebhookDeliveryHandler(
	webhookRepository repository.WebhookRepository,
	webhookDeliveryRepository repository.WebhookDeliveryRepository,
	sender webhook.Sender,
	maxAttempts int,
	retryDelay time.Duration,
) *SendWebhookDeliveryHandler {
	return &SendWebhookDeliveryHandler{
		webhookRepository:         webhookRepository,
		webhookDeliveryRepository: webhookDeliveryRepository,
		sender:                    sender,
		maxAttempts:               maxAttempts,
		retryDelay:                retryDelay,
	}
}

func (h *SendWebhookDeliveryHandler) Handle(ctx context.Context, cmd SendWebhookDelivery) (struct{}, error) {
	webhookDelivery, err := h.webhookDeliveryRepository.GetByID(ctx, cmd.WebhookDeliveryID)

	if err != nil {
		return struct{}{}, err
	}

	targetWebhook, err := h.webhookRepository.GetByID(ctx, webhookDelivery.WebhookID())

	if err != nil {
		return struct{}{}, err
	}

	responseStatus, sendErr := h.sender.Send(ctx, targetWebhook, webhookDelivery)

	if sendErr != nil {
		// Do not schedule another attempt while the application is shutting down.
		if ctx.Err() != nil {
			return struct{}{}, ctx.Err()
		}

		return struct{}{}, h.handleFailure(
			ctx,
			webhookDelivery,
			responseStatus,
			sendErr,
		)
	}

	if responseStatus == nil {
		return struct{}{}, errors.New("webhook sender returned neither response status nor error")
	}

	if deliverTransitionErr := webhookDelivery.Deliver(responseStatus); deliverTransitionErr != nil {
		return struct{}{}, deliverTransitionErr
	}

	if err := h.webhookDeliveryRepository.Update(ctx, webhookDelivery); err != nil {
		return struct{}{}, err
	}

	return struct{}{}, nil
}

func (h *SendWebhookDeliveryHandler) handleFailure(
	ctx context.Context,
	delivery *entity.WebhookDelivery,
	responseStatus *int,
	sendErr error,
) error {
	nextAttempt := delivery.Attempts() + 1

	if nextAttempt >= h.maxAttempts {
		if err := delivery.Fail(responseStatus, sendErr.Error()); err != nil {
			return errors.Join(sendErr, err)
		}
	} else {
		nextAttemptAt := time.Now().Add(time.Duration(nextAttempt) * h.retryDelay)

		if err := delivery.Retry(nextAttemptAt, responseStatus, sendErr.Error()); err != nil {
			return errors.Join(sendErr, err)
		}
	}

	if err := h.webhookDeliveryRepository.Update(ctx, delivery); err != nil {
		return errors.Join(sendErr, err)
	}

	// The failed HTTP attempt was successfully recorded as retry or final
	// failure, so the application command completed successfully.
	return nil
}
