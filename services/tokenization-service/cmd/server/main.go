package main

import (
    "fmt"
    "log"
    "tokenization-service/config"
    "tokenization-service/internal/crypto"
    "tokenization-service/internal/handler"
    "tokenization-service/internal/repository"
    "tokenization-service/internal/service"
    "tokenization-service/pkg/database"
    
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
    
    encryptor := crypto.NewEncryptor(cfg.EncryptionKey)
    tokenRepo := repository.NewTokenRepository(db)
    tokenService := service.NewTokenService(tokenRepo, encryptor)
    tokenHandler := handler.NewTokenHandler(tokenService)
    healthHandler := handler.NewHealthHandler(db)
    
    app := fiber.New(fiber.Config{
        AppName: "Tokenization Service v1.0",
    })
    
    app.Use(recover.New())
    app.Use(logger.New())
    app.Use(cors.New())
    
    health := app.Group("/health")
    health.Get("/liveness", healthHandler.Liveness)
    health.Get("/readiness", healthHandler.Readiness)
    
    api := app.Group("/api/tokens")
    api.Post("/", tokenHandler.Tokenize)
    api.Get("/:token", tokenHandler.GetToken)
    api.Post("/:token/detokenize", tokenHandler.Detokenize)
    
    log.Printf("Tokenization Service starting on port %s", cfg.ServerPort)
    if err := app.Listen("0.0.0.0:" + cfg.ServerPort); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
