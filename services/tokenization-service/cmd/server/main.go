package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "modernc.org/sqlite"

	"tokenization-service/internal/crypto"
	"tokenization-service/internal/handler"
	"tokenization-service/internal/repository"
	"tokenization-service/internal/service"
)

func main() {
	// Get config from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	keyVaultURI := os.Getenv("AZURE_KEY_VAULT_URI")
	keyName := os.Getenv("AZURE_KEY_NAME")

	// Initialize database
	db, err := sql.Open("sqlite", "./tokens.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Create schema
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tokens (
			id TEXT PRIMARY KEY,
			token TEXT NOT NULL,
			encrypted_pan BLOB NOT NULL,
			last4 TEXT NOT NULL,
			brand TEXT NOT NULL,
			expiry_month INTEGER NOT NULL,
			expiry_year INTEGER NOT NULL,
			merchant_id TEXT NOT NULL,
			created_at DATETIME NOT NULL
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

	// Create Fiber app
	app := fiber.New()

	// Register routes
	tokenHandler.RegisterRoutes(app)

	// Start server
	log.Printf("Server listening on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}