package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"ledger-service/internal/model"
	"ledger-service/internal/service"
	"log"

	"github.com/payments-sandbox/pkg/events"
)

type LedgerConsumer struct {
	svc *service.LedgerService
	bus events.EventBus
}

func NewLedgerConsumer(svc *service.LedgerService, bus events.EventBus) *LedgerConsumer {
	return &LedgerConsumer{
		svc: svc,
		bus: bus,
	}
}

func (c *LedgerConsumer) Start(ctx context.Context) {
	if err := c.bus.Subscribe(ctx, "payment.authorized", c.handlePaymentAuthorized); err != nil {
		log.Printf("Failed to subscribe to payment.authorized: %v", err)
	} else {
		log.Printf("Subscribed to payment.authorized")
	}
}

func (c *LedgerConsumer) handlePaymentAuthorized(ctx context.Context, event events.Event) error {
	log.Printf("Ledger processing payment: %s", event.ID)

	var payment struct {
		ID         string  `json:"id"`
		MerchantID string  `json:"merchant_id"`
		Amount     float64 `json:"amount"`
		Currency   string  `json:"currency"`
	}

	if err := json.Unmarshal(event.Data, &payment); err != nil {
		return fmt.Errorf("invalid payment data: %w", err)
	}

	// Create Double Entry
	// Debit: Platform Clearing (Asset)
	// Credit: Merchant Account (Liability)
	
	req := &model.CreateEntriesRequest{
		PaymentID:   payment.ID,
		Description: fmt.Sprintf("Payment authorized: %s", payment.ID),
		Entries: []model.EntryRequest{
			{
				AccountID: "platform-clearing",
				EntryType: model.DEBIT,
				Amount:    payment.Amount,
				Currency:  payment.Currency,
			},
			{
				AccountID: fmt.Sprintf("merchant-%s", payment.MerchantID),
				EntryType: model.CREDIT,
				Amount:    payment.Amount,
				Currency:  payment.Currency,
			},
		},
	}

	_, err := c.svc.CreateEntries(req)
	if err != nil {
		log.Printf("Failed to create ledger entries: %v", err)
		return err
	}

	log.Printf("Ledger entries created for payment %s", payment.ID)
	return nil
}
