package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/codememory1/d8r/internal/application/command"
	"github.com/codememory1/d8r/internal/domain/entity"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/internal/infrastructure/config"
	"github.com/codememory1/d8r/pkg/cqrs"
	"golang.org/x/sync/errgroup"
)

// InspectTaskWorker claims pending tasks and inspects them concurrently.
type InspectTaskWorker struct {
	logger             *slog.Logger
	config             *config.InspectionWorker
	taskClaimer        command.TaskClaimer
	inspectTaskHandler cqrs.CommandHandler[command.InspectTask, valueobject.ID]
}

// NewInspectTaskWorker creates a worker for processing pending task inspections.
func NewInspectTaskWorker(
	logger *slog.Logger,
	config *config.InspectionWorker,
	taskClaimer command.TaskClaimer,
	inspectTaskHandler cqrs.CommandHandler[command.InspectTask, valueobject.ID],
) *InspectTaskWorker {
	return &InspectTaskWorker{
		logger:             logger,
		config:             config,
		taskClaimer:        taskClaimer,
		inspectTaskHandler: inspectTaskHandler,
	}
}

// Run continuously claims and processes pending tasks until the context is
// canceled or an unrecoverable error occurs.
func (w *InspectTaskWorker) Run(ctx context.Context) error {
	for {
		tasks, err := w.taskClaimer.ClaimPending(ctx, int(w.config.Concurrency))

		if err != nil {
			return err
		}

		if len(tasks) == 0 {
			if err := w.wait(ctx); err != nil {
				return err
			}
		}

		if err := w.processTasks(ctx, tasks); err != nil {
			return err
		}
	}
}

// processTasks processes the claimed tasks concurrently up to the configured
// concurrency limit.
func (w *InspectTaskWorker) processTasks(ctx context.Context, tasks []*entity.Task) error {
	var group errgroup.Group

	group.SetLimit(int(w.config.Concurrency))

	for _, task := range tasks {
		group.Go(func() error {
			return w.processTask(ctx, task)
		})
	}

	return group.Wait()
}

// processTask inspects a task and transitions it to the failed state when the
// inspection cannot be completed.
func (w *InspectTaskWorker) processTask(ctx context.Context, task *entity.Task) error {
	_, err := w.inspectTaskHandler.Handle(ctx, command.InspectTask{
		TaskID: task.ID(),
	})

	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err != nil {
		w.logger.ErrorContext(
			ctx,
			"failed to inspect task",
			slog.String("task_id", task.ID().String()),
			slog.Any("error", err),
		)
	}

	return nil
}

// wait pauses polling until the polling interval elapses or the context is canceled.
func (w *InspectTaskWorker) wait(ctx context.Context) error {
	ticker := time.NewTimer(1 * time.Second)

	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ticker.C:
		return nil
	}
}
