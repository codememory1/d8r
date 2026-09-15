package restful

import (
	"errors"
	"net/http"

	"github.com/codememory1/d8r/pkg/restful/respond"
)

// HandlerFunc defines an HTTP handler that can return an error.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// HandlerAdapter converts error-returning handlers into standard HTTP handlers.
type HandlerAdapter struct {
	responder respond.Responder
}

// NewHandlerAdapter creates a handler adapter using the provided responder.
func NewHandlerAdapter(responder respond.Responder) *HandlerAdapter {
	return &HandlerAdapter{
		responder: responder,
	}
}

// Wrap converts HandlerFunc into http.HandlerFunc and handles returned errors.
func (a *HandlerAdapter) Wrap(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			a.writeError(w, err)
		}
	}
}

// writeError converts an application error into an HTTP error response.
func (a *HandlerAdapter) writeError(w http.ResponseWriter, err error) {
	if httpErr, ok := errors.AsType[Error](err); ok {
		_ = a.responder.Respond(w, httpErr.Status, respond.NewErrorBody(httpErr.Message))

		return
	}

	_ = a.responder.Respond(
		w,
		http.StatusInternalServerError,
		respond.NewErrorBody("Internal server error"),
	)
}
