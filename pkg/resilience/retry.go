package resilience

import (
	"context"
	"time"
)

// RetryFunc is the function to be retried.
type RetryFunc func() error

// Retry executes the given function with exponential backoff.
func Retry(ctx context.Context, maxAttempts int, baseDelay time.Duration, fn RetryFunc) error {
	var err error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if err = fn(); err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(baseDelay * (1 << uint(attempt))):
			// Continue to next attempt
		}
	}
	return err
}
