package resilience

import (
	"time"

	"github.com/sony/gobreaker"
)

// BreakerConfig holds the configuration for a circuit breaker.
type BreakerConfig struct {
	Name        string
	MaxRequests uint32
	Interval    time.Duration
	Timeout     time.Duration
	Threshold   uint32 // Number of failures to trigger open state
}

// NewCircuitBreaker creates a new gobreaker.CircuitBreaker with the given config.
func NewCircuitBreaker(cfg BreakerConfig) *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        cfg.Name,
		MaxRequests: cfg.MaxRequests,
		Interval:    cfg.Interval,
		Timeout:     cfg.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= cfg.Threshold
		},
	}

	return gobreaker.NewCircuitBreaker(settings)
}
