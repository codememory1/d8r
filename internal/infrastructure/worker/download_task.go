package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/codememory1/d8r/internal/application/command"
	"github.com/codememory1/d8r/pkg/cqrs"
)

// DownloadTaskWorker continuously polls and processes pending download tasks.
type DownloadTaskWorker struct {
	logger                     *slog.Logger
	downloadReadyTasksHandlers cqrs.CommandHandler[command.DownloadReadyTasks, struct{}]
	limit                      int
}

func NewDownloadTaskWorker(
	logger *slog.Logger,
	downloadReadyTasksHandlers cqrs.CommandHandler[command.DownloadReadyTasks, struct{}],
	limit int,
) *DownloadTaskWorker {
	return &DownloadTaskWorker{
		logger:                     logger,
		downloadReadyTasksHandlers: downloadReadyTasksHandlers,
		limit:                      limit,
	}
}

// Run claims pending tasks and processes them concurrently until the context
// is canceled or task processing returns an error.
func (w *DownloadTaskWorker) Run(ctx context.Context) error {
	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

	for {
		_, err := w.downloadReadyTasksHandlers.Handle(ctx, command.DownloadReadyTasks{
			Limit: w.limit,
		})

		if err != nil {
			w.logger.ErrorContext(ctx, "Failed to download ready tasks", slog.Any("error", err))
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
}
