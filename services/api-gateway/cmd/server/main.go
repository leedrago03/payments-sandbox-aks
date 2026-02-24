package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/payments-sandbox/pkg/logging"
	"github.com/payments-sandbox/pkg/resilience"
	"github.com/sony/gobreaker"
	"api-gateway/internal/middleware"
)

func main() {
	// Initialize Tracing
	logging.InitTracing()

	app := fiber.New()

	// Tracing Middleware: Extract trace from incoming request and pass to context
	app.Use(func(c *fiber.Ctx) error {
		ctx := logging.ExtractTraceFromFastHTTP(c.UserContext(), &c.Request().Header)
		c.SetUserContext(ctx)
		return c.Next()
	})

	app.Use(logger.New())

	// Basic health check
	app.Get("/health/liveness", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	app.Get("/health/readiness", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Service URLs (should be env vars in real deployment)
	tokenServiceURL := os.Getenv("TOKENIZATION_SERVICE_URL")
	if tokenServiceURL == "" {
		tokenServiceURL = "http://tokenization-service:3003"
	}

	paymentServiceURL := os.Getenv("PAYMENT_SERVICE_URL")
	if paymentServiceURL == "" {
		paymentServiceURL = "http://payment-service:8081"
	}

	merchantServiceURL := os.Getenv("MERCHANT_SERVICE_URL")
	if merchantServiceURL == "" {
		merchantServiceURL = "http://merchant-service:3002"
	}

	// Circuit Breakers
	tokenBreaker := resilience.NewCircuitBreaker(resilience.BreakerConfig{
		Name:      "token-service",
		Threshold: 3,
		Timeout:   30 * time.Second,
	})

	paymentBreaker := resilience.NewCircuitBreaker(resilience.BreakerConfig{
		Name:      "payment-service",
		Threshold: 3,
		Timeout:   30 * time.Second,
	})

	// Middleware
	authMiddleware := middleware.APIKeyAuth(merchantServiceURL)

	// Helper for resilient proxying
	proxyWithResilience := func(c *fiber.Ctx, url string, breaker *gobreaker.CircuitBreaker) error {
		_, err := breaker.Execute(func() (interface{}, error) {
			// Inject tracing headers before proxying
			logging.InjectTraceToFastHTTP(c.UserContext(), &c.Request().Header)

			if err := proxy.Do(c, url); err != nil {
				return nil, err
			}
			// Treat 5xx as failures for the circuit breaker
			if c.Response().StatusCode() >= 500 {
				return nil, fmt.Errorf("upstream service error: %d", c.Response().StatusCode())
			}
			return nil, nil
		})

		if err != nil {
			if err == gobreaker.ErrOpenState {
				return c.Status(fiber.StatusServiceUnavailable).SendString("Service Unavailable (Circuit Open)")
			}
			// For 500s that tripped the breaker, the response body is already set by proxy.Do
			// For connection errors, we return the error
			return err
		}
		return nil
	}

	// Routes
	// ... (rest of routes)
	app.Group("/v1/tokens", func(c *fiber.Ctx) error {
		return proxyWithResilience(c, tokenServiceURL+"/v1/tokens"+c.Params("*"), tokenBreaker)
	})
	
	app.Post("/v1/tokenize", authMiddleware, func(c *fiber.Ctx) error {
		return proxyWithResilience(c, tokenServiceURL+"/v1/tokenize", tokenBreaker)
	})

	// Public Merchant Routes (No Auth required for creation/login)
	app.All("/v1/merchants*", func(c *fiber.Ctx) error {
		path := c.Path()
		// Rewrite /v1/... to /api/...
		// e.g., /v1/merchants -> /api/merchants
		newPath := "/api" + path[3:]
		
		// Inject tracing headers
		logging.InjectTraceToFastHTTP(c.UserContext(), &c.Request().Header)
		
		return proxy.Do(c, merchantServiceURL+newPath)
	})

	app.All("/v1/payments*", authMiddleware, func(c *fiber.Ctx) error {
		path := c.Path() // /v1/payments...
		return proxyWithResilience(c, paymentServiceURL+path, paymentBreaker)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("API Gateway starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
