package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/codememory1/d8r/internal/application/command"
	"github.com/codememory1/d8r/internal/domain/entity"
	"github.com/codememory1/d8r/internal/infrastructure/config"
	"github.com/codememory1/d8r/pkg/cqrs"
	"golang.org/x/sync/errgroup"
)

// DownloadTaskWorker continuously polls and processes pending download tasks.
type DownloadTaskWorker struct {
	logger              *slog.Logger
	config              *config.DownloadWorker
	taskClaimer         command.TaskClaimer
	downloadTaskHandler cqrs.CommandHandler[command.DownloadTask, any]
}

func NewDownloadTaskWorker(
	logger *slog.Logger,
	config *config.DownloadWorker,
	taskClaimer command.TaskClaimer,
	downloadTaskHandler cqrs.CommandHandler[command.DownloadTask, any],
) *DownloadTaskWorker {
	return &DownloadTaskWorker{
		logger:              logger,
		config:              config,
		taskClaimer:         taskClaimer,
		downloadTaskHandler: downloadTaskHandler,
	}
}

// Run claims pending tasks and processes them concurrently until the context
// is canceled or task processing returns an error.
func (w *DownloadTaskWorker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			tasks, err := w.taskClaimer.ClaimReadyToDownload(ctx, int(w.config.Concurrency))

			if err != nil {
				return err
			}

			if len(tasks) == 0 {
				if err := w.wait(ctx); err != nil {
					return err
				}

				continue
			}

			if err := w.processTasks(ctx, tasks); err != nil {
				return err
			}
		}
	}
}

// processTasks concurrently processes a claimed batch of download tasks.
func (w *DownloadTaskWorker) processTasks(ctx context.Context, tasks []*entity.Task) error {
	var group errgroup.Group

	group.SetLimit(int(w.config.Concurrency))

	for _, task := range tasks {
		group.Go(func() error {
			return w.processTask(ctx, task)
		})
	}

	return group.Wait()
}

// processTask processes a single download task and logs task-specific failures.
func (w *DownloadTaskWorker) processTask(ctx context.Context, task *entity.Task) error {
	_, err := w.downloadTaskHandler.Handle(ctx, command.DownloadTask{
		TaskID: task.ID(),
	})

	// Context cancellation must stop the entire worker.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err != nil {
		w.logger.ErrorContext(
			ctx,
			"download task processing failed",
			slog.String("task_id", task.ID().String()),
			slog.Any("error", err),
		)
	}

	// A task-specific failure must not stop the worker.
	return nil
}

// wait pauses polling until the polling interval elapses or the context is canceled.
func (w *DownloadTaskWorker) wait(ctx context.Context) error {
	ticker := time.NewTimer(1 * time.Second)

	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ticker.C:
		return nil
	}
}
