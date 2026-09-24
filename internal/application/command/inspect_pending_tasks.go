package command

import (
	"context"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/cqrs"
	"golang.org/x/sync/errgroup"
)

var _ cqrs.CommandHandler[InspectPendingTasks, struct{}] = (*InspectPendingTasksHandler)(nil)

type InspectPendingTasks struct {
	Concurrency int
	Limit       int
}

type InspectPendingTasksHandler struct {
	taskClaimer        TaskClaimer
	inspectTaskHandler cqrs.CommandHandler[InspectTask, valueobject.ID]
}

func NewInspectPendingTasksHandler(
	taskClaimer TaskClaimer,
	inspectTaskHandler cqrs.CommandHandler[InspectTask, valueobject.ID],
) *InspectPendingTasksHandler {
	return &InspectPendingTasksHandler{
		taskClaimer:        taskClaimer,
		inspectTaskHandler: inspectTaskHandler,
	}
}

func (h *InspectPendingTasksHandler) Handle(ctx context.Context, cmd InspectPendingTasks) (struct{}, error) {
	taskIDs, err := h.taskClaimer.ClaimPending(ctx, cmd.Limit)

	if err != nil {
		return struct{}{}, err
	}

	var group errgroup.Group

	group.SetLimit(cmd.Concurrency)

	for _, taskID := range taskIDs {
		group.Go(func() error {
			_, err := h.inspectTaskHandler.Handle(ctx, InspectTask{
				TaskID: taskID,
			})

			return err
		})
	}

	return struct{}{}, group.Wait()
}
