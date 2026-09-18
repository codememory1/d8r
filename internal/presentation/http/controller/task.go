package controller

import (
	"errors"
	"net/http"

	"github.com/codememory1/d8r/internal/application/command"
	"github.com/codememory1/d8r/internal/application/query"
	"github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/internal/presentation/http/request"
	"github.com/codememory1/d8r/pkg/cqrs"
	"github.com/codememory1/d8r/pkg/restful"
	"github.com/codememory1/d8r/pkg/restful/respond"
	"github.com/go-chi/chi/v5"
)

// TaskController handles HTTP requests related to task management
// and delegates their execution to the corresponding CQRS handlers.
type TaskController struct {
	responder  respond.Responder
	createTask cqrs.CommandHandler[command.CreateTask, valueobject.ID]
	getTask    cqrs.QueryHandler[query.GetTask, query.GetTaskResult]
	listTasks  cqrs.QueryHandler[query.ListTasks, query.ListTasksResult]
}

// NewTaskController creates a task controller with its command and query handlers.
func NewTaskController(
	responder respond.Responder,
	createTask cqrs.CommandHandler[command.CreateTask, valueobject.ID],
	getTask cqrs.QueryHandler[query.GetTask, query.GetTaskResult],
	listTasks cqrs.QueryHandler[query.ListTasks, query.ListTasksResult],
) *TaskController {
	return &TaskController{
		responder:  responder,
		createTask: createTask,
		getTask:    getTask,
		listTasks:  listTasks,
	}
}

// Create handles a request to create a new httpdownload task.
func (c *TaskController) Create(w http.ResponseWriter, r *http.Request) error {
	req, err := restful.DecodeBody[request.CreateTask](r.Body)

	if err != nil {
		return err
	}

	cmd, err := req.ToCommand()

	if err != nil {
		return err
	}

	id, err := c.createTask.Handle(r.Context(), cmd)

	if err != nil {
		return err
	}

	return c.responder.Respond(w, http.StatusCreated, respond.NewSuccessBody(map[string]any{
		"id": id.String(),
	}))
}

// Get handles a request to retrieve a task by its identifier.
func (c *TaskController) Get(w http.ResponseWriter, r *http.Request) error {
	rawID := chi.URLParam(r, "id")
	id, err := valueobject.ParseID(rawID)

	if err != nil {
		return restful.NewError(http.StatusBadRequest, "invalid task id", err)
	}

	result, err := c.getTask.Handle(r.Context(), query.GetTask{
		ID: id,
	})

	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return restful.NewError(http.StatusNotFound, "task not found", err)
		}

		return err
	}

	return c.responder.Respond(w, http.StatusOK, respond.NewSuccessBody(result))
}

// List handles a request to retrieve a collection of tasks.
func (c *TaskController) List(w http.ResponseWriter, r *http.Request) error {
	req, err := restful.DecodeQuery[request.ListTasks](r.URL.Query())

	if err != nil {
		return err
	}

	q, err := req.ToQuery()

	if err != nil {
		return err
	}

	result, err := c.listTasks.Handle(r.Context(), q)

	if err != nil {
		return err
	}

	responseBody := respond.
		NewSuccessBody(result.Items).
		WithCursorPagination(q.Limit, result.NextCursor)

	return c.responder.Respond(w, http.StatusOK, responseBody)
}
