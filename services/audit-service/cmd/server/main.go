package main

import (
    "audit-service/config"
    "audit-service/internal/handler"
    "audit-service/internal/repository"
    "audit-service/internal/service"
    "audit-service/pkg/database"
    "fmt"
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
    
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
    
    db, err := database.InitDB(connStr)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Close()
    
    auditRepo := repository.NewAuditRepository(db)
    auditService := service.NewAuditService(auditRepo)
    auditHandler := handler.NewAuditHandler(auditService)
    healthHandler := handler.NewHealthHandler(db)
    
    app := fiber.New(fiber.Config{
        AppName: "Audit Service v1.0",
    })
    
    app.Use(recover.New())
    app.Use(logger.New())
    app.Use(cors.New())
    
    health := app.Group("/health")
    health.Get("/liveness", healthHandler.Liveness)
    health.Get("/readiness", healthHandler.Readiness)
    
    api := app.Group("/api/audit")
    api.Post("/logs", auditHandler.CreateLog)
    api.Get("/logs", auditHandler.QueryLogs)
    api.Get("/logs/entity/:entityId", auditHandler.GetEntityLogs)
    api.Get("/stats", auditHandler.GetStats)
    
    log.Printf("Audit Service starting on port %s", cfg.ServerPort)
    if err := app.Listen("0.0.0.0:" + cfg.ServerPort); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
