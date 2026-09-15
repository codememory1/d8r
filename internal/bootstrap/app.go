package bootstrap

import (
	"context"
	"log/slog"

	"github.com/codememory1/d8r/internal/application/command"
	"github.com/codememory1/d8r/internal/application/query"
	"github.com/codememory1/d8r/internal/infrastructure/config"
	"github.com/codememory1/d8r/internal/infrastructure/persistence/postgres"
	"github.com/codememory1/d8r/internal/infrastructure/persistence/postgres/repository"
	"github.com/codememory1/d8r/internal/presentation/http/controller"
	"github.com/codememory1/d8r/pkg/restful"
	"github.com/codememory1/d8r/pkg/restful/respond"
)

// Controllers contains the application's HTTP controllers.
type Controllers struct {
	Task *controller.TaskController
}

// App contains the application's configuration, infrastructure dependencies,
// and presentation-layer components.
type App struct {
	Config         config.Config
	Pool           *postgres.ConnectionPool
	Logger         *slog.Logger
	HandlerAdapter *restful.HandlerAdapter
	Controllers    Controllers
}

// NewApp initializes the application and composes its dependencies.
func NewApp(ctx context.Context, configuration config.Config, logger *slog.Logger) (*App, error) {
	postgresPool, err := postgres.NewConnectionPool(ctx, configuration.Postgres)

	if err != nil {
		return nil, err
	}

	app := &App{
		Config: configuration,
		Pool:   postgresPool,
		Logger: logger,
	}

	if err := app.compose(); err != nil {
		app.Close()

		return nil, err
	}

	return app, nil
}

// Close releases resources owned by the application.
func (a *App) Close() {
	if a.Pool != nil {
		a.Pool.Close()
	}
}

// compose creates and connects the application's dependencies.
func (a *App) compose() error {
	// Init Responders
	jsonResponder := respond.NewJSONResponder()

	// Init handler adapter
	a.HandlerAdapter = restful.NewHandlerAdapter(jsonResponder)

	// Init Repositories
	taskRepository := repository.NewTaskRepository(a.Pool)

	// Init Query/Command Handlers
	createTaskHandler := command.NewCreateTaskHandler(taskRepository)
	getTaskHandler := query.NewGetTaskHandler(taskRepository)
	listTasksHandler := query.NewListTasksHandler(taskRepository)

	// Init Controllers
	a.Controllers.Task = controller.NewTaskController(jsonResponder, createTaskHandler, getTaskHandler, listTasksHandler)

	return nil
}
