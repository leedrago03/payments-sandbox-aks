package main

import (
	"os"
	"payment-service/internal/handler"
	"payment-service/internal/integration"
	"payment-service/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/payments-sandbox/pkg/events"
	"github.com/payments-sandbox/pkg/logging"
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
    // ...
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
	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		eventBus = events.NewRedisBus(redisAddr, "", 0)
		logger.Info("Using Redis Event Bus", zap.String("addr", redisAddr))
	} else {
		eventBus = events.NewMemoryBus()
		logger.Info("Using In-Memory Event Bus")
	}

	paymentSvc := service.NewPaymentService(tokenClient, acquirerClient, eventBus, logger)
	
	app := fiber.New()
	
	h := handler.NewPaymentHandler(paymentSvc)
	h.RegisterRoutes(app)

	logger.Info("Payment Service starting", zap.String("port", port))
	if err := app.Listen(":" + port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
