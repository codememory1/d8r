package webhook

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	appwebhook "github.com/codememory1/d8r/internal/application/webhook"
	"github.com/codememory1/d8r/internal/domain/entity"
)

var _ appwebhook.Sender = (*HTTPSender)(nil)

// ErrUnexpectedResponseStatus indicates that the webhook endpoint returned
// a non-successful HTTP status.
var ErrUnexpectedResponseStatus = errors.New("unexpected webhook response status")

// HTTPSender sends webhook deliveries to their target endpoints over HTTP.
type HTTPSender struct {
	client *http.Client
}

// NewHTTPSender creates an HTTP-backed webhook sender.
func NewHTTPSender(client *http.Client) *HTTPSender {
	return &HTTPSender{
		client: client,
	}
}

// Send posts the delivery payload to the target webhook and returns the
// received HTTP status code.
func (s *HTTPSender) Send(ctx context.Context, webhook *entity.Webhook, delivery *entity.WebhookDelivery) (status *int, err error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		webhook.URL().String(),
		bytes.NewReader(delivery.Payload()),
	)

	if err != nil {
		return nil, fmt.Errorf("create webhook request: %w", err)
	}

	// Apply headers configured by the webhook owner.
	for name, value := range webhook.Headers().Map() {
		request.Header.Set(name, value)
	}

	// Set the default content type unless it was explicitly configured.
	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}

	// Add metadata that allows the receiver to identify and deduplicate
	// webhook deliveries.
	request.Header.Set("X-Webhook-Delivery-ID", delivery.ID().String())
	request.Header.Set("X-Webhook-Event", delivery.EventType().String())

	response, err := s.client.Do(request)

	if err != nil {
		return nil, fmt.Errorf("send webhook request: %w", err)
	}

	defer response.Body.Close()

	responseStatus := response.StatusCode

	if responseStatus < http.StatusOK || responseStatus >= http.StatusMultipleChoices {
		return &responseStatus, fmt.Errorf("%w: %d", ErrUnexpectedResponseStatus, responseStatus)
	}

	return &responseStatus, nil
}
