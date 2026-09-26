package bootstrap

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/codememory1/d8r/internal/application/command"
	appdownload "github.com/codememory1/d8r/internal/application/download"
	appevent "github.com/codememory1/d8r/internal/application/event"
	appinspect "github.com/codememory1/d8r/internal/application/inspect"
	"github.com/codememory1/d8r/internal/application/query"
	"github.com/codememory1/d8r/internal/application/storage"
	"github.com/codememory1/d8r/internal/application/transaction"
	domainevent "github.com/codememory1/d8r/internal/domain/event"
	domainrepo "github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/infrastructure/config"
	"github.com/codememory1/d8r/internal/infrastructure/eventbus"
	"github.com/codememory1/d8r/internal/infrastructure/eventcodec"
	"github.com/codememory1/d8r/internal/infrastructure/eventregistry"
	"github.com/codememory1/d8r/internal/infrastructure/httpdownload"
	"github.com/codememory1/d8r/internal/infrastructure/inspect"
	infraoutbox "github.com/codememory1/d8r/internal/infrastructure/outbox"
	"github.com/codememory1/d8r/internal/infrastructure/persistence/postgres"
	"github.com/codememory1/d8r/internal/infrastructure/persistence/postgres/outbox"
	"github.com/codememory1/d8r/internal/infrastructure/persistence/postgres/reader"
	postgresrepo "github.com/codememory1/d8r/internal/infrastructure/persistence/postgres/repository"
	"github.com/codememory1/d8r/internal/infrastructure/storage/filesystem"
	"github.com/codememory1/d8r/internal/infrastructure/worker"
	"github.com/codememory1/d8r/internal/presentation/http/controller"
	"github.com/codememory1/d8r/pkg/ddd"
	"github.com/codememory1/d8r/pkg/restful"
	"github.com/codememory1/d8r/pkg/restful/respond"
)

// Services contains shared application and infrastructure services.
type Services struct {
	Responder        respond.Responder
	Storage          storage.Storage
	EventDispatcher  appevent.Dispatcher
	EventPublisher   appevent.Publisher
	EventRegistry    infraoutbox.EventRegistry
	EventEncoder     infraoutbox.Encoder
	EventDecoder     infraoutbox.Decoder
	OutboxStore      infraoutbox.Store
	OutboxRelay      worker.OutboxRelay
	StrategySelector *appdownload.StrategySelector
	Inspector        appinspect.Inspector
	Downloader       appdownload.Downloader
}

// Repositories contains domain repository implementations used by the application.
type Repositories struct {
	Task           domainrepo.TaskRepository
	TaskInspection domainrepo.TaskInspectionRepository
	Webhook        domainrepo.WebhookRepository
}

// Claimers contains application dependencies responsible for claiming tasks
// that are ready for background processing.
type Claimers struct {
	PendingTask         command.PendingTaskClaimer
	ReadyToDownloadTask command.ReadyToDownloadTaskClaimer
}

// Readers contains read-side dependencies used by query handlers.
type Readers struct {
	Task query.TaskReader
}

// CommandHandlers contains handlers responsible for application commands.
type CommandHandlers struct {
	CreateTask         *command.CreateTaskHandler
	CreateWebhook      *command.CreateWebhookHandler
	InspectTask        *command.InspectTaskHandler
	InspectPendingTask *command.InspectPendingTasksHandler
	DownloadTask       *command.DownloadTaskHandler
	DownloadReadyTasks *command.DownloadReadyTasksHandler
}

// QueryHandlers contains handlers responsible for application queries.
type QueryHandlers struct {
	GetTask   *query.GetTaskHandler
	ListTasks *query.ListTasksHandler
}

// Controllers contains the application's HTTP controllers.
type Controllers struct {
	Task    *controller.TaskController
	Webhook *controller.WebhookController
}

// Workers contains background workers executed by the application.
type Workers struct {
	DownloadTask *worker.DownloadTaskWorker
	InspectTask  *worker.InspectTaskWorker
	OutboxEvent  *worker.OutboxEventWorker
}

// App is the application composition root and owns its dependencies.
type App struct {
	Config      config.Config
	Pool        *postgres.ConnectionPool
	Transaction transaction.Manager
	Logger      *slog.Logger

	Services     Services
	Repositories Repositories
	Claimers     Claimers
	Readers      Readers

	CommandHandlers CommandHandlers
	QueryHandlers   QueryHandlers

	HandlerAdapter *restful.HandlerAdapter
	Controllers    Controllers
	Workers        Workers
}

// NewApp initializes the application and composes its dependencies.
func NewApp(
	ctx context.Context,
	configuration config.Config,
	logger *slog.Logger,
) (*App, error) {
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

// compose initializes and connects application dependencies in dependency order.
func (a *App) compose() error {
	a.initServices()
	a.initPersistence()
	a.initReaders()

	a.initCommandHandlers()
	a.initQueryHandlers()

	a.initPresentation()
	a.initWorkers()

	return nil
}

// initServices initializes shared application and infrastructure services.
func (a *App) initServices() {
	eventRegistry := eventregistry.NewRegistry()

	eventRegistry.Register(domainevent.TaskCreatedType, func() ddd.Event {
		return &domainevent.TaskCreated{}
	})

	a.Services.Responder = respond.NewJSONResponder()
	a.Services.Storage = filesystem.NewStorage("./files")

	a.Services.EventRegistry = eventRegistry
	a.Services.EventEncoder = eventcodec.NewEncoder()
	a.Services.EventDecoder = eventcodec.NewDecoder()
	a.Services.EventDispatcher = eventbus.NewDispatcher()

	a.Services.EventPublisher = outbox.NewPublisher(
		a.Pool,
		a.Services.EventEncoder,
	)

	a.Services.OutboxStore = outbox.NewStore(a.Pool)

	a.Services.OutboxRelay = infraoutbox.NewRelay(
		a.Services.EventRegistry,
		a.Services.EventDecoder,
		a.Services.OutboxStore,
		a.Services.EventDispatcher,
		a.Config.Workers.OutboxEvent.Concurrency,
	)

	a.Services.StrategySelector = appdownload.NewStrategySelector(
		a.Config.Download.MinParallelSize.Bytes(),
	)

	a.Services.Inspector = inspect.NewHttpInspector(
		&http.Client{},
		inspect.Options{
			MinRangeProbeSize: a.Config.Download.MinParallelSize.Bytes(),
			ReservedHeaders:   make([]string, 0),
		},
	)

	a.Services.Downloader = httpdownload.NewHttpDownloader(
		&http.Client{},
		a.Services.Storage,
		httpdownload.Options{
			BufferSize:       a.Config.Download.BufferSize.Bytes(),
			MaxParallelParts: a.Config.Download.MaxParallelParts,
			RangeParts:       a.Config.Download.RangeParts,
		},
	)
}

// initPersistence initializes repository and task-claiming implementations.
func (a *App) initPersistence() {
	taskRepository := postgresrepo.NewTaskRepository(a.Pool)

	a.Repositories.Task = taskRepository
	a.Repositories.TaskInspection = postgresrepo.NewTaskInspectionRepository(a.Pool)
	a.Repositories.Webhook = postgresrepo.NewWebhookRepository(a.Pool)

	a.Claimers.PendingTask = taskRepository
	a.Claimers.ReadyToDownloadTask = taskRepository
}

// initReaders initializes read-side persistence dependencies.
func (a *App) initReaders() {
	a.Readers.Task = reader.NewTaskReader(a.Pool)
}

// initCommandHandlers initializes application command handlers.
func (a *App) initCommandHandlers() {
	a.CommandHandlers.CreateTask = command.NewCreateTaskHandler(
		a.Transaction,
		a.Services.EventPublisher,
		a.Repositories.Task,
	)

	a.CommandHandlers.InspectTask = command.NewInspectTaskHandler(
		a.Services.Inspector,
		a.Services.StrategySelector,
		a.Repositories.Task,
		a.Repositories.TaskInspection,
		a.Transaction,
	)

	a.CommandHandlers.DownloadTask = command.NewDownloadTaskHandler(
		a.Repositories.Task,
		a.Repositories.TaskInspection,
		a.Services.Downloader,
	)

	a.CommandHandlers.CreateWebhook = command.NewCreateWebhookHandler(
		a.Repositories.Webhook,
	)

	a.CommandHandlers.DownloadReadyTasks = command.NewDownloadReadyTasksHandler(
		a.Claimers.ReadyToDownloadTask,
		a.CommandHandlers.DownloadTask,
		a.Config.Workers.Download.Concurrency,
	)

	a.CommandHandlers.InspectPendingTask = command.NewInspectPendingTasksHandler(
		a.Claimers.PendingTask,
		a.CommandHandlers.InspectTask,
		a.Config.Workers.Inspection.Concurrency,
	)
}

// initQueryHandlers initializes application query handlers.
func (a *App) initQueryHandlers() {
	a.QueryHandlers.GetTask = query.NewGetTaskHandler(
		a.Readers.Task,
	)

	a.QueryHandlers.ListTasks = query.NewListTasksHandler(
		a.Readers.Task,
	)
}

// initPresentation initializes HTTP presentation-layer dependencies.
func (a *App) initPresentation() {
	a.HandlerAdapter = restful.NewHandlerAdapter(
		a.Services.Responder,
	)

	a.Controllers.Task = controller.NewTaskController(
		a.Services.Responder,
		a.CommandHandlers.CreateTask,
		a.QueryHandlers.GetTask,
		a.QueryHandlers.ListTasks,
	)

	a.Controllers.Webhook = controller.NewWebhookController(
		a.Services.Responder,
		a.CommandHandlers.CreateWebhook,
	)
}

// initWorkers initializes background application workers.
func (a *App) initWorkers() {
	a.Workers.DownloadTask = worker.NewDownloadTaskWorker(
		a.Logger,
		a.CommandHandlers.DownloadReadyTasks,
		a.Config.Workers.Download.Limit,
	)

	a.Workers.InspectTask = worker.NewInspectTaskWorker(
		a.Logger,
		a.CommandHandlers.InspectPendingTask,
		a.Config.Workers.Inspection.Limit,
	)

	a.Workers.OutboxEvent = worker.NewOutboxEventWorker(
		a.Logger,
		a.Services.OutboxRelay,
		a.Config.Workers.OutboxEvent.Limit,
	)
}
