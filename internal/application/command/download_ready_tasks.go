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
	Concurrency int
	Limit       int
}

type DownloadReadyTasksHandler struct {
	taskClaimer         ReadyToDownloadTaskClaimer
	downloadTaskHandler cqrs.CommandHandler[DownloadTask, struct{}]
}

func NewDownloadReadyTasksHandler(
	taskClaimer ReadyToDownloadTaskClaimer,
	downloadTaskHandler cqrs.CommandHandler[DownloadTask, struct{}],
) *DownloadReadyTasksHandler {
	return &DownloadReadyTasksHandler{
		taskClaimer:         taskClaimer,
		downloadTaskHandler: downloadTaskHandler,
	}
}

func (h *DownloadReadyTasksHandler) Handle(ctx context.Context, cmd DownloadReadyTasks) (struct{}, error) {
	taskIDs, err := h.taskClaimer.ClaimReadyToDownload(ctx, cmd.Limit)

	if err != nil {
		return struct{}{}, err
	}

	var group errgroup.Group

	group.SetLimit(cmd.Concurrency)

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
