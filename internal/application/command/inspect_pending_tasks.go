package command

import (
	"context"

	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/codememory1/d8r/pkg/cqrs"
	"golang.org/x/sync/errgroup"
)

var _ cqrs.CommandHandler[InspectPendingTasks, struct{}] = (*InspectPendingTasksHandler)(nil)

// PendingTaskClaimer defines a contract for atomically claiming tasks
// that are waiting for inspection.
type PendingTaskClaimer interface {
	ClaimPending(ctx context.Context, limit int) ([]valueobject.ID, error)
}

// InspectPendingTasks requests the claiming and inspection of a batch
// of pending tasks.
type InspectPendingTasks struct {
	Limit int
}

// InspectPendingTasksHandler claims pending tasks and inspects them
// concurrently.
type InspectPendingTasksHandler struct {
	taskClaimer        PendingTaskClaimer
	inspectTaskHandler cqrs.CommandHandler[InspectTask, valueobject.ID]
	concurrency        int
}

// NewInspectPendingTasksHandler creates a handler for processing batches
// of pending tasks.
func NewInspectPendingTasksHandler(
	taskClaimer PendingTaskClaimer,
	inspectTaskHandler cqrs.CommandHandler[InspectTask, valueobject.ID],
	concurrency int,
) *InspectPendingTasksHandler {
	return &InspectPendingTasksHandler{
		taskClaimer:        taskClaimer,
		inspectTaskHandler: inspectTaskHandler,
		concurrency:        concurrency,
	}
}

// Handle claims a batch of pending tasks and inspects them concurrently,
// respecting the configured concurrency limit.
func (h *InspectPendingTasksHandler) Handle(ctx context.Context, cmd InspectPendingTasks) (struct{}, error) {
	taskIDs, err := h.taskClaimer.ClaimPending(ctx, cmd.Limit)

	if err != nil {
		return struct{}{}, err
	}

	var group errgroup.Group

	group.SetLimit(h.concurrency)

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
