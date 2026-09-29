package command

import (
	"context"
	"errors"
	"time"

	"github.com/codememory1/d8r/internal/application/download"
	"github.com/codememory1/d8r/internal/application/event"
	"github.com/codememory1/d8r/internal/application/transaction"
	domainevent "github.com/codememory1/d8r/internal/domain/event"
	"github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/cqrs"
)

var _ cqrs.CommandHandler[DownloadTask, struct{}] = (*DownloadTaskHandler)(nil)

// DownloadTask represents a command to download the resource associated with a task.
type DownloadTask struct {
	TaskID valueobject.ID
}

// DownloadTaskHandler handles resource download commands and persists task state transitions.
type DownloadTaskHandler struct {
	tm                       transaction.Manager
	eventPublisher           event.Publisher
	taskRepository           repository.TaskRepository
	taskInspectionRepository repository.TaskInspectionRepository
	downloader               download.Downloader
}

// NewDownloadTaskHandler creates a handler for processing download tasks.
func NewDownloadTaskHandler(
	tm transaction.Manager,
	eventPublisher event.Publisher,
	taskRepository repository.TaskRepository,
	taskInspectionRepository repository.TaskInspectionRepository,
	downloader download.Downloader,
) *DownloadTaskHandler {
	return &DownloadTaskHandler{
		tm:                       tm,
		eventPublisher:           eventPublisher,
		taskRepository:           taskRepository,
		taskInspectionRepository: taskInspectionRepository,
		downloader:               downloader,
	}
}

// Handle downloads a resource using its latest inspection result.
func (h *DownloadTaskHandler) Handle(ctx context.Context, cmd DownloadTask) (struct{}, error) {
	// Load the task and its latest inspection result.
	task, err := h.taskRepository.GetByID(ctx, cmd.TaskID)

	if err != nil {
		return struct{}{}, err
	}

	taskInspection, err := h.taskInspectionRepository.GetLastByTaskID(ctx, cmd.TaskID)

	if err != nil {
		return struct{}{}, err
	}

	// Prefer the task's filename, falling back to the inspected filename.
	filename := task.Filename()

	if filename == nil {
		filename = taskInspection.Filename()
	}

	// Build download options from the task and inspection metadata.
	options := download.Options{
		URL:          taskInspection.EffectiveURL(),
		ContentType:  taskInspection.ContentType(),
		Headers:      task.Headers(),
		Filename:     filename,
		Size:         taskInspection.Size(),
		Strategy:     taskInspection.Strategy(),
		ETag:         taskInspection.ETag(),
		LastModified: taskInspection.LastModifiedAt(),
	}

	if downloadErr := h.downloader.Download(ctx, options); downloadErr != nil {
		// Propagate context cancellation without marking the task as failed.
		if ctx.Err() != nil {
			return struct{}{}, ctx.Err()
		}

		// Mark the task as failed and persist its state.
		// Preserve the download error if the transition or update also fails.
		if failTransitionErr := task.Fail(); failTransitionErr != nil {
			return struct{}{}, errors.Join(downloadErr, failTransitionErr)
		}

		// Persist the fail state and publish its event within the same transaction.
		transactionErr := h.tm.Run(ctx, func(ctx context.Context) error {
			if updateErr := h.taskRepository.Update(ctx, task); updateErr != nil {
				return errors.Join(downloadErr, updateErr)
			}

			return h.eventPublisher.Publish(ctx, domainevent.NewTaskDownloadFailed(
				task.ID(),
				time.Now(),
			))
		})

		if transactionErr != nil {
			return struct{}{}, errors.Join(downloadErr, transactionErr)
		}

		return struct{}{}, err
	}

	// Mark the task as completed after a successful download.
	if completeTransitionErr := task.Complete(); completeTransitionErr != nil {
		return struct{}{}, completeTransitionErr
	}

	// Persist the completed state and publish its event within the same transaction.
	transactionErr := h.tm.Run(ctx, func(ctx context.Context) error {
		if err := h.taskRepository.Update(ctx, task); err != nil {
			return err
		}

		return h.eventPublisher.Publish(ctx, domainevent.NewTaskDownloadCompleted(
			task.ID(),
			time.Now(),
		))
	})

	if transactionErr != nil {
		return struct{}{}, transactionErr
	}

	return struct{}{}, nil
}
