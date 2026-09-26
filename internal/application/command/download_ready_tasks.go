package command

import (
	"context"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/cqrs"
	"golang.org/x/sync/errgroup"
)

var _ cqrs.CommandHandler[DownloadReadyTasks, struct{}] = (*DownloadReadyTasksHandler)(nil)

type ReadyToDownloadTaskClaimer interface {
	ClaimReadyToDownload(ctx context.Context, limit int) ([]valueobject.ID, error)
}

type DownloadReadyTasks struct {
	Limit int
}

type DownloadReadyTasksHandler struct {
	taskClaimer         ReadyToDownloadTaskClaimer
	downloadTaskHandler cqrs.CommandHandler[DownloadTask, struct{}]
	concurrency         int
}

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
