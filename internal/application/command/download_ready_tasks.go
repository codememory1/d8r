package command

import (
	"context"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/cqrs"
	"golang.org/x/sync/errgroup"
)

var _ cqrs.CommandHandler[DownloadReadyTasks, struct{}] = (*DownloadReadyTasksHandler)(nil)

// ReadyToDownloadTaskClaimer defines a contract for atomically claiming tasks
// that are ready to be downloaded.
type ReadyToDownloadTaskClaimer interface {
	ClaimReadyToDownload(ctx context.Context, limit int) ([]valueobject.ID, error)
}

// DownloadReadyTasks requests the claiming and processing of a batch of tasks
// that are ready to be downloaded.
type DownloadReadyTasks struct {
	Limit int
}

// DownloadReadyTasksHandler claims ready tasks and processes their downloads
// concurrently.
type DownloadReadyTasksHandler struct {
	taskClaimer         ReadyToDownloadTaskClaimer
	downloadTaskHandler cqrs.CommandHandler[DownloadTask, struct{}]
	concurrency         int
}

// NewDownloadReadyTasksHandler creates a handler for processing batches of
// tasks that are ready to be downloaded.
func NewDownloadReadyTasksHandler(
	taskClaimer ReadyToDownloadTaskClaimer,
	downloadTaskHandler cqrs.CommandHandler[DownloadTask, struct{}],
	concurrency int,
) *DownloadReadyTasksHandler {
	return &DownloadReadyTasksHandler{
		taskClaimer:         taskClaimer,
		downloadTaskHandler: downloadTaskHandler,
		concurrency:         concurrency,
	}
}

// Handle claims a batch of ready tasks and downloads them concurrently,
// respecting the configured concurrency limit.
func (h *DownloadReadyTasksHandler) Handle(ctx context.Context, cmd DownloadReadyTasks) (struct{}, error) {
	taskIDs, err := h.taskClaimer.ClaimReadyToDownload(ctx, cmd.Limit)

	if err != nil {
		return struct{}{}, err
	}

	var group errgroup.Group

	group.SetLimit(h.concurrency)

	for _, taskID := range taskIDs {
		group.Go(func() error {
			_, err := h.downloadTaskHandler.Handle(ctx, DownloadTask{
				TaskID: taskID,
			})

			return err
		})
	}

	return struct{}{}, group.Wait()
}
