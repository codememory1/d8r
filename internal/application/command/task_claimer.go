package command

import (
	"context"

	"github.com/codememory1/d8r/internal/domain/valueobject"
)

// TaskClaimer atomically claims tasks for background processing.
type TaskClaimer interface {
	// ClaimPending claims tasks waiting for inspection.
	ClaimPending(ctx context.Context, limit int) ([]valueobject.ID, error)

	// ClaimReadyToDownload claims tasks waiting for download.
	ClaimReadyToDownload(ctx context.Context, limit int) ([]valueobject.ID, error)
}
