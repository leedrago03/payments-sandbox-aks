package main

import (
	"crypto/tls"
	"os"
	"payment-service/internal/handler"
	"payment-service/internal/integration"
	"payment-service/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/payments-sandbox/pkg/events"
	"github.com/payments-sandbox/pkg/logging"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	logger, err := logging.NewLogger("payment-service")
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	
	// Service URLs (should be env vars in real deployment)
	tokenSvcURL := os.Getenv("TOKENIZATION_SERVICE_URL")
	if tokenSvcURL == "" {
		tokenSvcURL = "http://localhost:3003"
	}

	acquirerSvcURL := os.Getenv("ACQUIRER_SERVICE_URL")
	if acquirerSvcURL == "" {
		acquirerSvcURL = "http://localhost:3004"
	}

	// Initialize dependencies
	tokenClient := integration.NewTokenizationClient(tokenSvcURL)
	acquirerClient := integration.NewAcquirerClient(acquirerSvcURL)
	
	var eventBus events.EventBus
	var redisClient *redis.Client

	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		redisPassword := os.Getenv("REDIS_PASSWORD")
		useTLS := os.Getenv("REDIS_TLS") == "true"
		
		// Event Bus
		eventBus = events.NewRedisBus(redisAddr, redisPassword, 0, useTLS)
		logger.Info("Using Redis Event Bus", zap.String("addr", redisAddr), zap.Bool("tls", useTLS))

		// Redis Client for Idempotency
		opt := &redis.Options{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       0,
		}
		if useTLS {
			opt.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
			}
		}
		redisClient = redis.NewClient(opt)

	} else {
		eventBus = events.NewMemoryBus()
		logger.Info("Using In-Memory Event Bus")
		// For local dev without Redis, we might want a mock or nil (idempotency won't work)
		logger.Warn("Redis not configured - Idempotency will be disabled")
	}

	paymentSvc := service.NewPaymentService(tokenClient, acquirerClient, eventBus, logger, redisClient)
	
	app := fiber.New()
	
	h := handler.NewPaymentHandler(paymentSvc)
	h.RegisterRoutes(app)

	// Add health checks
	app.Get("/health/liveness", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	app.Get("/health/readiness", func(c *fiber.Ctx) error {
		// In production, you'd check DB/Redis connectivity here
		return c.SendStatus(fiber.StatusOK)
	})

	logger.Info("Payment Service starting", zap.String("port", port))
	if err := app.Listen(":" + port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
