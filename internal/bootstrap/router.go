package bootstrap

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Router creates and returns an HTTP handler with the application's routes registered.
func (a *App) Router() http.Handler {
	r := chi.NewRouter()

	r.Post("/tasks", a.HandlerAdapter.Wrap(a.Controllers.Task.Create))
	r.Get("/tasks/{id}", a.HandlerAdapter.Wrap(a.Controllers.Task.Get))
	r.Get("/tasks", a.HandlerAdapter.Wrap(a.Controllers.Task.List))

	return r
}
