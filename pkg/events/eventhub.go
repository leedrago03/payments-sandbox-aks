package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
)

// EventHubBus is an Azure Event Hubs implementation of EventBus.
type EventHubBus struct {
	producerClient *azeventhubs.ProducerClient
	consumerClient *azeventhubs.ConsumerClient
	namespace      string
	hubName        string
}

// NewEventHubBus creates a new EventHubBus using Azure Identity (Workload Identity).
// Note: In a real Event Hubs setup, a single client usually connects to a single Hub.
// Our interface assumes a generic "Publish" where type might map to Hub or Topic.
// For simplicity in this sandbox, we will assume one "events" Hub for all messages,
// or we would need a registry of clients.
// Strategy: Use a single Hub named 'payment-events' for all traffic, filter by property.
func NewEventHubBus(namespace, hubName string) (*EventHubBus, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create azure credential: %w", err)
	}

	// Create Producer
	producer, err := azeventhubs.NewProducerClient(namespace, hubName, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	// Create Consumer (using default consumer group)
	// Note: Creating consumer here is just a placeholder. Real consumption requires
	// a Processor and Blob Storage for checkpoints.
	// For this Step 1 implementation, we will focus on PRODUCING first, 
	// as Consumption in Event Hubs is complex (ProcessorClient) and requires storage account.
	
	return &EventHubBus{
		producerClient: producer,
		namespace:      namespace,
		hubName:        hubName,
	}, nil
}

// Publish sends an event to the Event Hub.
func (b *EventHubBus) Publish(ctx context.Context, event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	batch, err := b.producerClient.NewEventDataBatch(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create batch: %w", err)
	}

	eventData := &azeventhubs.EventData{
		Body: payload,
		Properties: map[string]any{
			"EventType": event.Type,
		},
	}

	if err := batch.AddEventData(eventData, nil); err != nil {
		return fmt.Errorf("failed to add event to batch: %w", err)
	}

	if err := b.producerClient.SendEventDataBatch(ctx, batch, nil); err != nil {
		return fmt.Errorf("failed to send batch: %w", err)
	}

	return nil
}

// Subscribe is a stub for now. 
// Azure Event Hubs consumption typically requires a ProcessorClient with Blob Storage for checkpoints.
// Implementing a full Processor here is significant work involving Azure Storage SDK.
// For Step 1, we acknowledge this complexity.
func (b *EventHubBus) Subscribe(ctx context.Context, eventType string, handler Handler) error {
	return fmt.Errorf("subscription via EventHubBus not fully implemented in Step 1. Requires Checkpoint Store")
}

// Close closes the clients.
func (b *EventHubBus) Close() error {
	return b.producerClient.Close(context.Background())
}
