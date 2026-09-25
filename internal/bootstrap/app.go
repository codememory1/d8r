package bootstrap

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/codememory1/d8r/internal/application/command"
	"github.com/codememory1/d8r/internal/application/download"
	"github.com/codememory1/d8r/internal/application/query"
	"github.com/codememory1/d8r/internal/application/transaction"
	"github.com/codememory1/d8r/internal/infrastructure/config"
	"github.com/codememory1/d8r/internal/infrastructure/eventbus"
	"github.com/codememory1/d8r/internal/infrastructure/eventcodec"
	"github.com/codememory1/d8r/internal/infrastructure/httpdownload"
	"github.com/codememory1/d8r/internal/infrastructure/inspect"
	infraoutbox "github.com/codememory1/d8r/internal/infrastructure/outbox"
	"github.com/codememory1/d8r/internal/infrastructure/persistence/postgres"
	"github.com/codememory1/d8r/internal/infrastructure/persistence/postgres/outbox"
	"github.com/codememory1/d8r/internal/infrastructure/persistence/postgres/reader"
	"github.com/codememory1/d8r/internal/infrastructure/persistence/postgres/repository"
	"github.com/codememory1/d8r/internal/infrastructure/storage/filesystem"
	"github.com/codememory1/d8r/internal/presentation/http/controller"
	"github.com/codememory1/d8r/internal/presentation/worker"
	"github.com/codememory1/d8r/pkg/restful"
	"github.com/codememory1/d8r/pkg/restful/respond"
)

// Controllers contains the application's HTTP controllers.
type Controllers struct {
	Task    *controller.TaskController
	Webhook *controller.WebhookController
}

type Workers struct {
	DownloadTask *worker.DownloadTaskWorker
	InspectTask  *worker.InspectTaskWorker
	OutboxEvent  *worker.OutboxEventWorker
}

// App contains the application's configuration, infrastructure dependencies,
// and presentation-layer components.
type App struct {
	Config         config.Config
	Pool           *postgres.ConnectionPool
	Transaction    transaction.Manager
	Logger         *slog.Logger
	HandlerAdapter *restful.HandlerAdapter
	Controllers    Controllers
	Workers        Workers
}

// NewApp initializes the application and composes its dependencies.
func NewApp(ctx context.Context, configuration config.Config, logger *slog.Logger) (*App, error) {
	postgresPool, err := postgres.NewConnectionPool(ctx, configuration.Postgres)

	if err != nil {
		return nil, err
	}

	app := &App{
		Config:      configuration,
		Pool:        postgresPool,
		Transaction: postgresPool,
		Logger:      logger,
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
	taskInspectionRepository := repository.NewTaskInspectionRepository(a.Pool)
	webhookRepository := repository.NewWebhookRepository(a.Pool)

	// Init Readers
	taskReader := reader.NewTaskReader(a.Pool)

	// Init clients
	client := http.Client{}

	// Init Storages
	fsStorage := filesystem.NewStorage("./files")

	// Init Downloader
	httpDownloader := httpdownload.NewHttpDownloader(&client, fsStorage, &a.Config.Download)

	// Init Inspector
	httpInspector := inspect.NewHttpInspector(&client, &a.Config.Download, make([]string, 0))

	// Init Event Dispatcher
	eventDispatcher := eventbus.NewDispatcher()

	// Init Event Decoder
	eventDecoder := eventcodec.NewDecoder()

	// Event Publisher
	eventPublisher := outbox.NewPublisher(a.Pool)

	// Outbox
	outboxStore := outbox.NewStore(a.Pool)
	outboxRelay := infraoutbox.NewRelay(outboxStore, eventDecoder, eventDispatcher, 10)

	// Init strategy selector
	strategySelector := download.NewStrategySelector(
		a.Config.Download.MinParallelSize.Bytes(),
	)

	// Init Query/Command Handlers
	createTaskHandler := command.NewCreateTaskHandler(a.Pool, eventPublisher, taskRepository)
	getTaskHandler := query.NewGetTaskHandler(taskReader)
	listTasksHandler := query.NewListTasksHandler(taskReader)
	inspectTaskHandler := command.NewInspectTaskHandler(
		httpInspector,
		strategySelector,
		taskRepository,
		taskInspectionRepository,
		a.Transaction,
	)
	downloadTaskHandler := command.NewDownloadTaskHandler(taskRepository, taskInspectionRepository, httpDownloader)
	createWebhookHandler := command.NewCreateWebhookHandler(webhookRepository)
	downloadReadyTasksHandlers := command.NewDownloadReadyTasksHandler(taskRepository, downloadTaskHandler)
	inspectPendingTasksHandler := command.NewInspectPendingTasksHandler(taskRepository, inspectTaskHandler)

	// Init Controllers
	a.Controllers.Task = controller.NewTaskController(jsonResponder, createTaskHandler, getTaskHandler, listTasksHandler)
	a.Controllers.Webhook = controller.NewWebhookController(jsonResponder, createWebhookHandler)

	// Init Workers
	a.Workers.DownloadTask = worker.NewDownloadTaskWorker(
		a.Logger,
		downloadReadyTasksHandlers,
		a.Config.Workers.Download.Concurrency,
		a.Config.Workers.Download.Limit,
	)
	a.Workers.InspectTask = worker.NewInspectTaskWorker(
		a.Logger,
		inspectPendingTasksHandler,
		a.Config.Workers.Inspection.Concurrency,
		a.Config.Workers.Inspection.Limit,
	)
	a.Workers.OutboxEvent = worker.NewOutboxEventWorker(a.Pool, a.Logger, outboxRelay)

	return nil
}
