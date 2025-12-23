package events

import (
	"context"
)

// Handler is a function that processes an event.
type Handler func(ctx context.Context, event Event) error

// EventBus defines the interface for publishing and subscribing to events.
type EventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(ctx context.Context, eventType string, handler Handler) error
	Close() error
}
