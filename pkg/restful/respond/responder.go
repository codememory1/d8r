package respond

import "net/http"

// Responder writes API responses to an HTTP response writer.
type Responder interface {
	Respond(w http.ResponseWriter, status int, body Body) error
}
