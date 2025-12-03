package main

import (
    "acquirer-simulator/config"
    "acquirer-simulator/internal/handler"
    "acquirer-simulator/internal/service"
    "log"
    
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    acquirerService := service.NewAcquirerService(cfg.SuccessRate, cfg.TimeoutRate)
    acquirerHandler := handler.NewAcquirerHandler(acquirerService)
    
    app := fiber.New(fiber.Config{
        AppName: "Acquirer Simulator v1.0",
    })
    
    app.Use(recover.New())
    app.Use(logger.New())
    app.Use(cors.New())
    
    app.Get("/health", acquirerHandler.Health)
    
    api := app.Group("/api/acquirer")
    api.Post("/authorize", acquirerHandler.Authorize)
    api.Post("/capture", acquirerHandler.Capture)
    api.Post("/refund", acquirerHandler.Refund)
    
    log.Printf("Acquirer Simulator starting on port %s", cfg.ServerPort)
    if err := app.Listen("0.0.0.0:" + cfg.ServerPort); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
