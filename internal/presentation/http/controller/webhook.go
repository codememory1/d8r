package controller

import (
	"net/http"

	"github.com/codememory1/d8r/internal/application/command"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/internal/presentation/http/request"
	"github.com/codememory1/d8r/pkg/cqrs"
	"github.com/codememory1/d8r/pkg/restful"
	"github.com/codememory1/d8r/pkg/restful/respond"
)

// WebhookController handles HTTP requests related to webhook management.
type WebhookController struct {
	responder            respond.Responder
	createWebhookHandler cqrs.CommandHandler[command.CreateWebhook, valueobject.ID]
}

// NewWebhookController creates a controller for webhook-related HTTP requests.
func NewWebhookController(
	responder respond.Responder,
	createWebhookHandler cqrs.CommandHandler[command.CreateWebhook, valueobject.ID],
) *WebhookController {
	return &WebhookController{
		responder:            responder,
		createWebhookHandler: createWebhookHandler,
	}
}

// Create handles a request to create a new webhook.
func (c *WebhookController) Create(w http.ResponseWriter, r *http.Request) error {
	req, err := restful.DecodeBody[request.CreateWebhook](r.Body)

	if err != nil {
		return err
	}

	cmd, err := req.ToCommand()

	if err != nil {
		return err
	}

	id, err := c.createWebhookHandler.Handle(r.Context(), cmd)

	if err != nil {
		return err
	}

	return c.responder.Respond(w, http.StatusCreated, respond.NewSuccessBody(map[string]any{
		"id": id.String(),
	}))
}
