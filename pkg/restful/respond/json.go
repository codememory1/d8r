package respond

import (
	"encoding/json"
	"net/http"
)

// JSONResponder writes API responses encoded as JSON.
type JSONResponder struct{}

// NewJSONResponder creates a JSON API responder.
func NewJSONResponder() JSONResponder {
	return JSONResponder{}
}

// Respond writes the response body as JSON with the provided HTTP status.
func (r JSONResponder) Respond(w http.ResponseWriter, status int, body Body) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(body)
}
