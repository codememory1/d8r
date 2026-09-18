package request

import (
	"errors"
	"fmt"

	"github.com/codememory1/d8r/internal/application/command"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/restful"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Ensure CreateWebhookRequest implements the expected HTTP input contract.
var _ restful.Input = (*CreateWebhook)(nil)

// CreateWebhook represents the request body for creating a webhook.
type CreateWebhook struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Events  []string          `json:"events"`
}

// Validate validates the webhook URL and subscribed event types.
func (r CreateWebhook) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.URL, validation.Required.Error("URL is a required field.")),
		validation.Field(&r.Events,
			validation.Required.Error("At least one event is required."),
			validation.Each(validation.Required.Error("Event type cannot be empty.")),
			validation.By(validateUniqueEvents),
		),
	)
}

// validateUniqueEvents ensures that each event type occurs only once.
func validateUniqueEvents(value any) error {
	events, ok := value.([]string)

	if !ok {
		return errors.New("events must be an array of strings")
	}

	seen := make(map[string]struct{}, len(events))

	for _, event := range events {
		if _, exists := seen[event]; exists {
			return fmt.Errorf("event %q is duplicated", event)
		}

		seen[event] = struct{}{}
	}

	return nil
}

// ToCommand converts the validated request into a CreateWebhook command.
func (r CreateWebhook) ToCommand() (command.CreateWebhook, error) {
	url, err := valueobject.NewURL(r.URL)

	if err != nil {
		return command.CreateWebhook{}, err
	}

	headers, err := valueobject.NewHeaders(r.Headers)

	if err != nil {
		return command.CreateWebhook{}, err
	}

	// Convert raw event names into validated domain value objects.
	eventTypes := make([]valueobject.WebhookEventType, len(r.Events))

	for i, event := range r.Events {
		eventType, err := valueobject.NewWebhookEventType(event)

		if err != nil {
			return command.CreateWebhook{}, err
		}

		eventTypes[i] = eventType
	}

	return command.CreateWebhook{
		URL:        url,
		Headers:    headers,
		EventTypes: eventTypes,
	}, nil
}
