package service

import (
	"context"
	"fmt"
	"payment-service/internal/integration"
	"payment-service/internal/model"
	"time"

	"github.com/google/uuid"
	"github.com/payments-sandbox/pkg/events"
	"github.com/payments-sandbox/pkg/resilience"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

type PaymentService struct {
	tokenClient    *integration.TokenizationClient
	acquirerClient *integration.AcquirerClient
	eventBus       events.EventBus
	breaker        *gobreaker.CircuitBreaker
	logger         *zap.Logger
}

func NewPaymentService(tokenClient *integration.TokenizationClient, acquirerClient *integration.AcquirerClient, eventBus events.EventBus, logger *zap.Logger) *PaymentService {
	cb := resilience.NewCircuitBreaker(resilience.BreakerConfig{
		Name:      "acquirer-breaker",
		Threshold: 3, // Fail after 3 consecutive errors
		Timeout:   30 * time.Second,
	})

	return &PaymentService{
		tokenClient:    tokenClient,
		acquirerClient: acquirerClient,
		eventBus:       eventBus,
		breaker:        cb,
		logger:         logger,
	}
}

func (s *PaymentService) Authorize(ctx context.Context, req *model.CreatePaymentRequest) (*model.Payment, error) {
	// 1. Tokenize if PAN is provided (mocking this flow for now if token already exists)
	var token string
	if req.Token != "" {
		token = req.Token
	} else {
		// Call Tokenization service
		// In a real implementation, we would extract PAN from request (which shouldn't be here in prod flow usually)
		// For this sandbox, we assume the frontend might send a token, or we tokenize here.
		// Let's assume we call tokenization if we had PAN.
		token = "tok_mock_" + uuid.New().String()
	}

	// 2. Call Acquirer
	result, err := s.breaker.Execute(func() (interface{}, error) {
		return s.acquirerClient.Authorize(integration.AuthRequest{
			Amount:     req.Amount,
			Currency:   req.Currency,
			CardToken:  token,
			MerchantID: req.MerchantID,
		})
	})
	
	var authResp *integration.AuthResponse
	if err == nil {
		authResp = result.(*integration.AuthResponse)
	}
	
	status := model.StatusAuthorized
	if err != nil {
		status = model.StatusFailed
		s.logger.Error("Acquirer call failed", 
			zap.Error(err), 
			zap.String("merchant_id", req.MerchantID),
			zap.Float64("amount", req.Amount))
	} else if authResp.Status != "APPROVED" {
		status = model.StatusFailed
		s.logger.Warn("Acquirer declined payment", 
			zap.String("merchant_id", req.MerchantID),
			zap.String("status", authResp.Status))
	} else {
		s.logger.Info("Payment authorized successfully", 
			zap.String("payment_id", authResp.TransactionID),
			zap.String("merchant_id", req.MerchantID))
	}

	// 3. Create Payment record
	payment := &model.Payment{
		ID:             uuid.New().String(),
		MerchantID:     req.MerchantID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Token:          token,
		Status:         status,
		IdempotencyKey: req.IdempotencyKey,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// 4. Publish Event
	event, _ := events.NewEvent("payment.authorized", "payment-service", payment)
	if status == model.StatusFailed {
		event, _ = events.NewEvent("payment.failed", "payment-service", payment)
	}
	
	s.eventBus.Publish(ctx, *event)

	if status == model.StatusFailed {
		if err != nil {
			return nil, err
		}
		// Return payment with failed status
		return payment, nil
	}

	return payment, nil
}
