package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/codememory1/d8r/internal/application/command"
	"github.com/codememory1/d8r/pkg/cqrs"
)

// InspectTaskWorker claims pending tasks and inspects them concurrently.
type InspectTaskWorker struct {
	logger                     *slog.Logger
	inspectPendingTasksHandler cqrs.CommandHandler[command.InspectPendingTasks, struct{}]
	limit                      int
}

// NewInspectTaskWorker creates a worker for processing pending task inspections.
func NewInspectTaskWorker(
	logger *slog.Logger,
	inspectPendingTasksHandler cqrs.CommandHandler[command.InspectPendingTasks, struct{}],
	limit int,
) *InspectTaskWorker {
	return &InspectTaskWorker{
		logger:                     logger,
		inspectPendingTasksHandler: inspectPendingTasksHandler,
		limit:                      limit,
	}
}

// Run continuously claims and processes pending tasks until the context is
// canceled or an unrecoverable error occurs.
func (w *InspectTaskWorker) Run(ctx context.Context) error {
	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

	for {
		_, err := w.inspectPendingTasksHandler.Handle(ctx, command.InspectPendingTasks{
			Limit: w.limit,
		})

		if err != nil {
			w.logger.ErrorContext(ctx, "Failed to inspect pending tasks", slog.Any("error", err))
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
}
