package main

import (
    "context"
    "fmt"
    "ledger-service/config"
    "ledger-service/internal/consumer"
    "ledger-service/internal/handler"
    "ledger-service/internal/repository"
    "ledger-service/internal/service"
    "ledger-service/pkg/database"
    "log"
    "os"
    
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/recover"
    "github.com/payments-sandbox/pkg/events"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
    
    db, err := database.InitDB(connStr)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Close()
    
    // Initialize Event Bus
    var eventBus events.EventBus
    if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
        eventBus = events.NewRedisBus(redisAddr, "", 0)
        log.Printf("Using Redis Event Bus at %s", redisAddr)
    } else {
        eventBus = events.NewMemoryBus()
        log.Println("Using In-Memory Event Bus")
    }
    defer eventBus.Close()

    ledgerRepo := repository.NewLedgerRepository(db)
    ledgerService := service.NewLedgerService(ledgerRepo)
    ledgerHandler := handler.NewLedgerHandler(ledgerService)
    healthHandler := handler.NewHealthHandler(db)
    
    // Start Consumer
    ledgerConsumer := consumer.NewLedgerConsumer(ledgerService, eventBus)
    ledgerConsumer.Start(context.Background())

    app := fiber.New(fiber.Config{
        AppName: "Ledger Service v1.0",
    })
    
    app.Use(recover.New())
    app.Use(logger.New())
    app.Use(cors.New())
    
    health := app.Group("/health")
    health.Get("/liveness", healthHandler.Liveness)
    health.Get("/readiness", healthHandler.Readiness)
    
    api := app.Group("/api/ledger")
    api.Post("/entries", ledgerHandler.CreateEntries)
    api.Get("/payments/:paymentId/entries", ledgerHandler.GetPaymentEntries)
    api.Get("/accounts/:accountId/balance", ledgerHandler.GetAccountBalance)
    api.Get("/balances", ledgerHandler.GetAllBalances)
    
    log.Printf("Ledger Service starting on port %s", cfg.ServerPort)
    if err := app.Listen("0.0.0.0:" + cfg.ServerPort); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
