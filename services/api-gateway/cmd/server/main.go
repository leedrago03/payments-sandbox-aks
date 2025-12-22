package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())

	// Basic health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Service URLs (should be env vars in real deployment)
	tokenServiceURL := os.Getenv("TOKEN_SERVICE_URL")
	if tokenServiceURL == "" {
		tokenServiceURL = "http://tokenization-service:3003"
	}

	paymentServiceURL := os.Getenv("PAYMENT_SERVICE_URL")
	if paymentServiceURL == "" {
		paymentServiceURL = "http://payment-service:8080"
	}

	// Routes
	// Tokenization
	app.Group("/v1/tokens", func(c *fiber.Ctx) error {
		return proxy.Do(c, tokenServiceURL+"/v1/tokens"+c.Params("*"))
	})
	
	// Direct proxying for now. In a real scenario, we'd map paths more explicitly.
	// For example:
	app.Post("/v1/tokenize", func(c *fiber.Ctx) error {
		return proxy.Do(c, tokenServiceURL+"/v1/tokenize")
	})

	app.All("/v1/payments*", func(c *fiber.Ctx) error {
		path := c.Path() // /v1/payments...
		return proxy.Do(c, paymentServiceURL+path)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("API Gateway starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
