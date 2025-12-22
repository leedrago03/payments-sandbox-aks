package main

import (
	"log"
	"os"
	"payment-service/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := fiber.New()
	
	h := handler.NewPaymentHandler()
	h.RegisterRoutes(app)

	log.Printf("Payment Service starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
