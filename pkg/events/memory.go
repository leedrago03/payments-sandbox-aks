package events

import (
	"context"
	"sync"
)

// MemoryBus is an in-memory implementation of EventBus.
type MemoryBus struct {
	mu          sync.RWMutex
	subscribers map[string][]Handler
}

// NewMemoryBus creates a new MemoryBus.
func NewMemoryBus() *MemoryBus {
	return &MemoryBus{
		subscribers: make(map[string][]Handler),
	}
}

// Publish publishes an event to all subscribers.
func (b *MemoryBus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	handlers, ok := b.subscribers[event.Type]
	b.mu.RUnlock()

	if !ok {
		return nil
	}

	for _, handler := range handlers {
		// In a real scenario, we might want to run this in a goroutine
		// or handle errors differently.
		if err := handler(ctx, event); err != nil {
			// Log error but continue with other handlers?
			// For now, just return it.
			return err
		}
	}

	return nil
}

// Subscribe adds a handler for an event type.
func (b *MemoryBus) Subscribe(ctx context.Context, eventType string, handler Handler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers[eventType] = append(b.subscribers[eventType], handler)
	return nil
}

// Close closes the bus.
func (b *MemoryBus) Close() error {
	return nil
}
