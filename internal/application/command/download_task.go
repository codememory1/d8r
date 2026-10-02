package command

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/codememory1/d8r/internal/application/download"
	"github.com/codememory1/d8r/internal/application/event"
	"github.com/codememory1/d8r/internal/application/transaction"
	"github.com/codememory1/d8r/internal/domain/entity"
	domainevent "github.com/codememory1/d8r/internal/domain/event"
	"github.com/codememory1/d8r/internal/domain/repository"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/cqrs"
	"github.com/codememory1/d8r/pkg/optional"
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

	// Build lifecycle callbacks for the download.
	lifecycle, finishLifecycle := h.buildLifecycle(ctx, cmd.TaskID, *taskInspection)
	defer finishLifecycle()

	downloadErr := h.downloader.Download(ctx, options, lifecycle)

	// Stop progress publishing and wait for the final callback to finish.
	finishLifecycle()

	if downloadErr != nil {
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

		return struct{}{}, downloadErr
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

// buildLifecycle builds download lifecycle callbacks and starts periodic progress publishing.
// It returns the lifecycle and a function that stops publishing and waits for the worker to finish.
func (h *DownloadTaskHandler) buildLifecycle(ctx context.Context, taskID valueobject.ID, inspection entity.TaskInspection) (download.Lifecycle, func()) {
	// Track counters updated concurrently by download callbacks.
	var downloadedBytes atomic.Int64
	var activeRequests atomic.Int64

	// Keep the last successfully published values.
	// These variables are accessed only by the publishing worker.
	var oldDownloadedBytes int64
	var oldActiveRequests int64

	stopPublishing := h.startPeriodic(ctx, 1, func() {
		// Capture the current counters for comparison and publishing.
		currentDownloadedBytes := downloadedBytes.Load()
		currentActiveRequests := activeRequests.Load()

		// Skip event publication if there have been no changes since the last publication.
		if currentDownloadedBytes != oldDownloadedBytes || currentActiveRequests != oldActiveRequests {
			// Publish the captured progress snapshot.
			publishErr := h.eventPublisher.Publish(ctx, domainevent.NewTaskDownloadProgress(
				taskID,
				currentDownloadedBytes,
				optional.Map(inspection.Size(), valueobject.ByteSize.Int64),
				activeRequests.Load(),
			))

			// Update the previous values only after successful publishing.
			// On failure, leave them unchanged so a later tick can retry.
			if publishErr == nil {
				oldDownloadedBytes = currentDownloadedBytes
				oldActiveRequests = currentActiveRequests
			}
		}
	})

	return download.Lifecycle{
		OnProgress: func(writtenBytes int64) {
			downloadedBytes.Add(writtenBytes)
		},
		OnActiveRequestsChanged: func(delta int64) {
			activeRequests.Add(delta)
		},
	}, stopPublishing
}

// startPeriodic runs onTick periodically until the context is canceled or a stop signal is received.
// It returns a function that signals the worker to stop and waits for it to finish.
func (h *DownloadTaskHandler) startPeriodic(ctx context.Context, interval time.Duration, onTick func()) func() {
	var wg sync.WaitGroup
	var stopOnce sync.Once

	// Closing this channel signals the worker to stop.
	stopCh := make(chan struct{})

	// Register the worker before starting its goroutine.
	wg.Add(1)

	go func() {
		// Notify waiting callers when the worker exits.
		defer wg.Done()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-stopCh:
				onTick()

				return
			case <-ticker.C:
				onTick()
			}
		}
	}()

	return func() {
		stopOnce.Do(func() {
			close(stopCh)
		})

		wg.Wait()
	}
}
