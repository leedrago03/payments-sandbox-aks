package service

import (
	"context"
	"encoding/json"
	"errors"
	"payment-service/internal/integration"
	"payment-service/internal/model"
	"time"

	"github.com/google/uuid"
	"github.com/payments-sandbox/pkg/events"
	"github.com/payments-sandbox/pkg/resilience"
	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

type PaymentService struct {
	tokenClient    *integration.TokenizationClient
	acquirerClient *integration.AcquirerClient
	eventBus       events.EventBus
	breaker        *gobreaker.CircuitBreaker
	logger         *zap.Logger
	redisClient    *redis.Client
}

func NewPaymentService(tokenClient *integration.TokenizationClient, acquirerClient *integration.AcquirerClient, eventBus events.EventBus, logger *zap.Logger, redisClient *redis.Client) *PaymentService {
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
		redisClient:    redisClient,
	}
}

func (s *PaymentService) Authorize(ctx context.Context, req *model.CreatePaymentRequest) (*model.Payment, error) {
	// 0. Idempotency Check
	idempotencyKey := req.IdempotencyKey
	if idempotencyKey != "" {
		key := "idempotency:" + idempotencyKey
		val, err := s.redisClient.Get(ctx, key).Result()
		if err == nil {
			if val == "PROCESSING" {
				return nil, errors.New("request already in progress")
			}
			// Return cached response
			var cachedPayment model.Payment
			if err := json.Unmarshal([]byte(val), &cachedPayment); err == nil {
				s.logger.Info("Idempotency hit", zap.String("key", idempotencyKey))
				return &cachedPayment, nil
			}
		} else if err != redis.Nil {
			s.logger.Error("Redis error checking idempotency", zap.Error(err))
			// Fail open or closed? Closed for safety in payments.
			return nil, errors.New("internal system error")
		}

		// Lock
		success, err := s.redisClient.SetNX(ctx, key, "PROCESSING", 5*time.Minute).Result()
		if err != nil {
			return nil, err
		}
		if !success {
			return nil, errors.New("request already in progress")
		}
		
		// Ensure we clear processing state if we crash/panic (though strict consistency might prefer letting it timeout)
		// Defer clearing is tricky because we want to overwrite with success. 
		// We'll handle overwrite at the end.
	}

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
	
	if err := s.eventBus.Publish(ctx, *event); err != nil {
		s.logger.Error("Failed to publish event to bus", 
			zap.Error(err), 
			zap.String("event_type", event.Type))
	} else {
		s.logger.Info("Event published successfully", zap.String("event_id", event.ID))
	}

	// 5. Save Idempotency Result
	if idempotencyKey != "" {
		jsonPayment, _ := json.Marshal(payment)
		s.redisClient.Set(ctx, "idempotency:"+idempotencyKey, jsonPayment, 24*time.Hour)
	}

	if status == model.StatusFailed {
		if err != nil {
			return nil, err
		}
		// Return payment with failed status
		return payment, nil
	}

	return payment, nil
}
