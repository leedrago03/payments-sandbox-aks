package main

import (
    "fmt"
    "log"
    "reconciliation-service/config"
    "reconciliation-service/internal/handler"
    "reconciliation-service/internal/repository"
    "reconciliation-service/internal/service"
    "reconciliation-service/pkg/database"
    
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
    
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
    
    db, err := database.InitDB(connStr)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Close()
    
    reconRepo := repository.NewReconciliationRepository(db)
    reconService := service.NewReconciliationService(reconRepo)
    reconHandler := handler.NewReconciliationHandler(reconService)
    healthHandler := handler.NewHealthHandler(db)
    
    app := fiber.New(fiber.Config{
        AppName: "Reconciliation Service v1.0",
    })
    
    app.Use(recover.New())
    app.Use(logger.New())
    app.Use(cors.New())
    
    health := app.Group("/health")
    health.Get("/liveness", healthHandler.Liveness)
    health.Get("/readiness", healthHandler.Readiness)
    
    api := app.Group("/api/reconciliation")
    api.Post("/reports", reconHandler.CreateReconciliation)
    api.Get("/reports/:reportId", reconHandler.GetReport)
    api.Get("/reports", reconHandler.GetAllReports)
    
    log.Printf("Reconciliation Service starting on port %s", cfg.ServerPort)
    if err := app.Listen("0.0.0.0:" + cfg.ServerPort); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
