package command

import (
	"context"
	"errors"

	"github.com/codememory1/d8r/internal/application/download"
	"github.com/codememory1/d8r/internal/application/inspect"
	"github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/cqrs"
)

var _ cqrs.CommandHandler[DownloadTask, any] = (*DownloadTaskHandler)(nil)

// DownloadTask represents a command to download the resource associated with a task.
type DownloadTask struct {
	TaskID valueobject.ID
}

// DownloadTaskHandler handles resource download commands and persists task state transitions.
type DownloadTaskHandler struct {
	taskRepository           repository.TaskRepository
	taskInspectionRepository repository.TaskInspectionRepository
	downloader               download.Downloader
}

// NewDownloadTaskHandler creates a handler for processing download tasks.
func NewDownloadTaskHandler(
	taskRepository repository.TaskRepository,
	taskInspectionRepository repository.TaskInspectionRepository,
	downloader download.Downloader,
) *DownloadTaskHandler {
	return &DownloadTaskHandler{
		taskRepository:           taskRepository,
		taskInspectionRepository: taskInspectionRepository,
		downloader:               downloader,
	}
}

// Handle downloads a resource using its latest inspection result.
func (h *DownloadTaskHandler) Handle(ctx context.Context, cmd DownloadTask) (any, error) {
	task, err := h.taskRepository.GetByID(ctx, cmd.TaskID)

	if err != nil {
		return nil, err
	}

	taskInspection, err := h.taskInspectionRepository.GetLastByTaskID(ctx, cmd.TaskID)

	if err != nil {
		return nil, err
	}

	options := download.Options{
		InspectionResult: inspect.Result{
			EffectiveURL:     taskInspection.EffectiveURL(),
			ContentType:      taskInspection.ContentType(),
			Filename:         taskInspection.Filename(),
			Size:             taskInspection.Size(),
			DownloadStrategy: taskInspection.Strategy(),
			ETag:             taskInspection.ETag(),
			LastModified:     taskInspection.LastModifiedAt(),
		},
		Headers:  task.Headers().Map(),
		Filename: task.Filename(),
	}

	if downloadErr := h.downloader.Download(ctx, options); downloadErr != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		if failTransitionErr := task.Fail(); failTransitionErr != nil {
			return nil, errors.Join(downloadErr, failTransitionErr)
		}

		if updateErr := h.taskRepository.Update(ctx, task); updateErr != nil {
			return nil, errors.Join(downloadErr, updateErr)
		}

		return nil, err
	}

	if completeTransitionErr := task.Complete(); completeTransitionErr != nil {
		return nil, completeTransitionErr
	}

	if err := h.taskRepository.Update(ctx, task); err != nil {
		return nil, err
	}

	return nil, nil
}
