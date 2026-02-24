package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"tokenization-service/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"

	"github.com/payments-sandbox/pkg/crypto"
	"github.com/payments-sandbox/pkg/logging"
	"tokenization-service/internal/handler"
	"tokenization-service/internal/repository"
	"tokenization-service/internal/service"
)

func main() {
	// Initialize Tracing
	logging.InitTracing()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	keyVaultURI := os.Getenv("AZURE_KEY_VAULT_URI")
	keyName := os.Getenv("AZURE_KEY_NAME")

	// Initialize database
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Create schema
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tokens (
			id TEXT PRIMARY KEY,
			token TEXT NOT NULL,
			encrypted_pan BYTEA NOT NULL,
			last4 TEXT NOT NULL,
			brand TEXT NOT NULL,
			expiry_month INTEGER NOT NULL,
			expiry_year INTEGER NOT NULL,
			merchant_id TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL
		);
	`)
	if err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}


	// Initialize crypto service
	encryptor, err := crypto.NewEncryptor(keyVaultURI, keyName)
	if err != nil {
		log.Fatalf("failed to create encryptor: %v", err)
	}

	// Initialize repository
	repo := repository.NewTokenRepository(db)

	// Initialize service
	tokenService := service.NewTokenService(repo, encryptor)

	// Initialize handler
	tokenHandler := handler.NewTokenHandler(tokenService)
	healthHandler := handler.NewHealthHandler(db)

	// Create Fiber app
	app := fiber.New()

	app.Use(recover.New())
	
	// Tracing Middleware: Extract trace from incoming request and pass to context
	app.Use(func(c *fiber.Ctx) error {
		ctx := logging.ExtractTraceFromFastHTTP(c.UserContext(), &c.Request().Header)
		c.SetUserContext(ctx)
		return c.Next()
	})

	app.Use(logger.New())

	// Register routes
	tokenHandler.RegisterRoutes(app)

	health := app.Group("/health")
	health.Get("/liveness", healthHandler.Liveness)
	health.Get("/readiness", healthHandler.Readiness)

	// Start server
	log.Printf("Server listening on port %s", cfg.ServerPort)
	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}