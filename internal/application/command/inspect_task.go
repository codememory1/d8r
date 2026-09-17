package command

import (
	"context"

	"github.com/codememory1/d8r/internal/application/inspect"
	"github.com/codememory1/d8r/internal/application/transaction"
	"github.com/codememory1/d8r/internal/domain/entity"
	"github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/cqrs"
)

// Compile-time check that InspectResourceHandler implements the required command handler.
var _ cqrs.CommandHandler[InspectTask, valueobject.ID] = (*InspectTaskHandler)(nil)

// InspectTask requests an inspection of the resource associated with a http download task.
type InspectTask struct {
	TaskID valueobject.ID
}

// InspectTaskHandler coordinates resource inspection and persists its result.
type InspectTaskHandler struct {
	inspector                inspect.Inspector
	taskRepository           repository.TaskRepository
	taskInspectionRepository repository.TaskInspectionRepository
	transactionManager       transaction.Manager
}

// NewInspectTaskHandler creates a handler for the InspectResource command.
func NewInspectTaskHandler(
	inspector inspect.Inspector,
	taskRepository repository.TaskRepository,
	taskInspectionRepository repository.TaskInspectionRepository,
	transactionManager transaction.Manager,
) *InspectTaskHandler {
	return &InspectTaskHandler{
		inspector:                inspector,
		taskRepository:           taskRepository,
		taskInspectionRepository: taskInspectionRepository,
		transactionManager:       transactionManager,
	}
}

// Handle inspects the task resource, stores the resulting inspection,
// and returns the identifier of the created inspection.
func (h *InspectTaskHandler) Handle(ctx context.Context, cmd InspectTask) (valueobject.ID, error) {
	// Load the task to obtain the original resource URL and related settings.
	task, err := h.taskRepository.GetById(ctx, cmd.TaskID)

	if err != nil {
		return valueobject.ID{}, err
	}

	// Inspect the remote resource and determine its http download metadata and strategy.
	result, err := h.inspector.Inspect(ctx, task.URL().String(), task.Headers().Map())

	if err != nil {
		return valueobject.ID{}, err
	}

	// Convert the application-level inspection result into a domain entity.
	taskInspection := entity.NewTaskInspection(
		cmd.TaskID,
		result.EffectiveURL,
		result.ContentType,
		result.Filename,
		result.Size,
		result.DownloadStrategy,
		result.ETag,
		result.LastModified,
	)

	task.ReadyToDownload()

	err = h.transactionManager.Run(ctx, func(ctx context.Context) error {
		if updateErr := h.taskRepository.Update(ctx, task); updateErr != nil {
			return updateErr
		}

		// Persist the inspection so it can be reused without repeated HTTP requests.
		if err := h.taskInspectionRepository.Save(ctx, taskInspection); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return valueobject.ID{}, err
	}

	return taskInspection.ID(), nil
}
