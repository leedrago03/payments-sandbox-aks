package main

import (
    "log"
    "merchant-service/config"
    "merchant-service/internal/handler"
    "merchant-service/internal/repository"
    "merchant-service/internal/service"
    "merchant-service/pkg/database"
    
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Initialize database
    db, err := database.InitDB(cfg.GetDBConnectionString())
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Close()
    
    // Initialize repository, service, and handler
    merchantRepo := repository.NewMerchantRepository(db)
    merchantService := service.NewMerchantService(merchantRepo)
    merchantHandler := handler.NewMerchantHandler(merchantService)
    healthHandler := handler.NewHealthHandler(db)
    
    // Create Fiber app
    app := fiber.New(fiber.Config{
        AppName: "Merchant Service v1.0",
    })
    
    // Middleware
    app.Use(recover.New())
    app.Use(logger.New())
    app.Use(cors.New())
    
    // Health check routes
    health := app.Group("/health")
    health.Get("/liveness", healthHandler.Liveness)
    health.Get("/readiness", healthHandler.Readiness)
    
    // API routes
    api := app.Group("/api/merchants")
    api.Post("/", merchantHandler.CreateMerchant)
    api.Get("/:id", merchantHandler.GetMerchant)
    api.Put("/:id", merchantHandler.UpdateMerchant)
    api.Post("/:id/api-keys", merchantHandler.CreateAPIKey)
    api.Get("/:id/api-keys", merchantHandler.GetAPIKeys)
    
    // Internal routes (for Gateway/other services)
    app.Post("/internal/api-keys/verify", merchantHandler.VerifyAPIKey)

    // Start server
    log.Printf("Merchant Service starting on port %s", cfg.ServerPort)
    if err := app.Listen("0.0.0.0:" + cfg.ServerPort); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
