package httpdownload

import (
	"context"
	"time"
)

// Retry executes operations repeatedly when they return an error.
type Retry struct {
	MaxAttempts int
	Delay       time.Duration
}

// Do execute operation until it succeeds or all attempts are exhausted.
func (r Retry) Do(ctx context.Context, operation func() error) error {
	lastErr := operation()

	if lastErr != nil {
		for i := 0; i < r.MaxAttempts; i++ {
			err := operation()

			if err == nil {
				return nil
			}

			lastErr = err

			timer := time.NewTimer(r.Delay)

			select {
			case <-ctx.Done():
				timer.Stop()

				return ctx.Err()

			case <-timer.C:
			}
		}
	}

	return lastErr
}
