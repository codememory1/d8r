package command

import (
	"context"
	"fmt"
	"time"

	"github.com/codememory1/d8r/internal/application/event"
	"github.com/codememory1/d8r/internal/application/transaction"
	domainevent "github.com/codememory1/d8r/internal/domain/event"
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
	tm                 transaction.Manager
	taskClaimer        PendingTaskClaimer
	inspectTaskHandler cqrs.CommandHandler[InspectTask, valueobject.ID]
	eventPublisher     event.Publisher
	concurrency        int
}

// NewInspectPendingTasksHandler creates a handler for processing batches
// of pending tasks.
func NewInspectPendingTasksHandler(
	tm transaction.Manager,
	taskClaimer PendingTaskClaimer,
	inspectTaskHandler cqrs.CommandHandler[InspectTask, valueobject.ID],
	eventPublisher event.Publisher,
	concurrency int,
) *InspectPendingTasksHandler {
	return &InspectPendingTasksHandler{
		tm:                 tm,
		taskClaimer:        taskClaimer,
		inspectTaskHandler: inspectTaskHandler,
		eventPublisher:     eventPublisher,
		concurrency:        concurrency,
	}
}

// Handle claims a batch of pending tasks and inspects them concurrently,
// respecting the configured concurrency limit.
func (h *InspectPendingTasksHandler) Handle(ctx context.Context, cmd InspectPendingTasks) (struct{}, error) {
	var taskIDs []valueobject.ID

	err := h.tm.Run(ctx, func(ctx context.Context) error {
		var err error

		taskIDs, err = h.taskClaimer.ClaimPending(ctx, cmd.Limit)

		if err != nil {
			return err
		}

		for _, taskID := range taskIDs {
			publishErr := h.eventPublisher.Publish(ctx, domainevent.NewTaskInspectionStarted(
				taskID,
				time.Now(),
			), new(taskID.String()))

			if publishErr != nil {
				return fmt.Errorf("failed event publication for task: %s: %w", taskID, publishErr)
			}
		}

		return nil
	})

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
