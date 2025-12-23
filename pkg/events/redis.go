package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RedisBus is a Redis-backed implementation of EventBus.
type RedisBus struct {
	client *redis.Client
}

// NewRedisBus creates a new RedisBus.
func NewRedisBus(addr string, password string, db int) *RedisBus {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisBus{client: client}
}

// Publish publishes an event to a Redis channel (topic).
// We use the event Type as the channel name.
func (b *RedisBus) Publish(ctx context.Context, event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return b.client.Publish(ctx, event.Type, payload).Err()
}

// Subscribe subscribes to a Redis channel (event type).
// It starts a goroutine to listen for messages.
func (b *RedisBus) Subscribe(ctx context.Context, eventType string, handler Handler) error {
	pubsub := b.client.Subscribe(ctx, eventType)
	
	// Check if subscription was successful
	if _, err := pubsub.Receive(ctx); err != nil {
		return err
	}

	go func() {
		defer pubsub.Close()
		ch := pubsub.Channel()

		for msg := range ch {
			var event Event
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				// TODO: Log error properly
				fmt.Printf("Error unmarshalling event: %v\n", err)
				continue
			}

			// Execute handler with background context (or new context) 
			// as the original subscription context might be long-lived or cancelled separately.
			// Ideally we pass a context that controls the lifecycle of the consumer.
			if err := handler(context.Background(), event); err != nil {
				fmt.Printf("Error handling event %s: %v\n", event.ID, err)
			}
		}
	}()

	return nil
}

// Close closes the Redis client.
func (b *RedisBus) Close() error {
	return b.client.Close()
}
