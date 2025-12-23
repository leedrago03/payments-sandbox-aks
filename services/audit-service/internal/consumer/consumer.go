package consumer

import (
	"audit-service/internal/model"
	"audit-service/internal/service"
	"context"
	"encoding/json"
	"log"

	"github.com/payments-sandbox/pkg/events"
)

type AuditConsumer struct {
	svc *service.AuditService
	bus events.EventBus
}

func NewAuditConsumer(svc *service.AuditService, bus events.EventBus) *AuditConsumer {
	return &AuditConsumer{
		svc: svc,
		bus: bus,
	}
}

func (c *AuditConsumer) Start(ctx context.Context) {
	// List of events to listen to
	eventTypes := []string{
		"payment.authorized",
		"payment.captured",
		"payment.failed",
		// Add more as needed
	}

	for _, eventType := range eventTypes {
		if err := c.bus.Subscribe(ctx, eventType, c.handleEvent); err != nil {
			log.Printf("Failed to subscribe to %s: %v", eventType, err)
		} else {
			log.Printf("Subscribed to %s", eventType)
		}
	}
}

func (c *AuditConsumer) handleEvent(ctx context.Context, event events.Event) error {
	log.Printf("Received event: %s (%s)", event.Type, event.ID)

	// Map event to Audit Log
	// We might want to extract specific details based on event type
	// For now, we dump the raw event data into 'Details'
	
	req := &model.CreateAuditLogRequest{
		EventType:  model.EventType(event.Type), // Assuming mapping matches or we cast
		EntityType: "EVENT", // Or derive from event source
		EntityID:   event.ID,
		ActorType:  "SYSTEM", // Or extract from metadata
		ActorID:    event.Source,
		Action:     "EVENT_RECEIVED",
		Details:    string(event.Data), // Store the payload
		Success:    true,
	}

	// If it's a payment event, we might want to be more specific
	if event.Type == "payment.authorized" {
		req.EntityType = "PAYMENT"
		// Parse payload to get Payment ID?
		var paymentData map[string]interface{}
		if err := json.Unmarshal(event.Data, &paymentData); err == nil {
			if id, ok := paymentData["id"].(string); ok {
				req.EntityID = id
			}
		}
	}

	_, err := c.svc.LogEvent(req)
	if err != nil {
		log.Printf("Failed to log audit event: %v", err)
		return err
	}

	return nil
}
