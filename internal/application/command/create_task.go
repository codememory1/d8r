package command

import (
	"context"

	appevent "github.com/codememory1/d8r/internal/application/event"
	"github.com/codememory1/d8r/internal/application/transaction"
	"github.com/codememory1/d8r/internal/domain/entity"
	"github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/cqrs"
)

// Ensure CreateTaskHandler implements the expected command handler contract.
var _ cqrs.CommandHandler[CreateTask, valueobject.ID] = (*CreateTaskHandler)(nil)

// CreateTask contains the validated values required to create a task.
type CreateTask struct {
	URL      valueobject.URL
	Headers  valueobject.Headers
	Filename *valueobject.Filename
	Priority valueobject.Priority
}

// CreateTaskHandler handles task creation commands.
type CreateTaskHandler struct {
	tm             transaction.Manager
	eventPublisher appevent.Publisher
	taskRepository repository.TaskRepository
}

// NewCreateTaskHandler creates a handler for task creation commands.
func NewCreateTaskHandler(
	tm transaction.Manager,
	eventPublisher appevent.Publisher,
	taskRepository repository.TaskRepository,
) *CreateTaskHandler {
	return &CreateTaskHandler{
		tm:             tm,
		eventPublisher: eventPublisher,
		taskRepository: taskRepository,
	}
}

// Handle creates a new task and persists it in the repository.
func (h *CreateTaskHandler) Handle(ctx context.Context, cmd CreateTask) (valueobject.ID, error) {
	task := entity.NewTask(
		cmd.URL,
		cmd.Headers,
		cmd.Filename,
		cmd.Priority,
	)
	events := task.PullEvents()

	err := h.tm.Run(ctx, func(ctx context.Context) error {
		if err := h.taskRepository.Save(ctx, task); err != nil {
			return err
		}

		for _, event := range events {
			if err := h.eventPublisher.Publish(ctx, event); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return valueobject.ID{}, err
	}

	return task.ID(), nil
}
