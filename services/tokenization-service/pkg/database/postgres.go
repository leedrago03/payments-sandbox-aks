package database

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/lib/pq"
)

func InitDB(connStr string) (*sql.DB, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    log.Println("Database connected successfully")
    
    if err := createTables(db); err != nil {
        return nil, fmt.Errorf("failed to create tables: %w", err)
    }
    
    return db, nil
}

func createTables(db *sql.DB) error {
    schema := `
    CREATE TABLE IF NOT EXISTS tokens (
        id VARCHAR(36) PRIMARY KEY,
        token VARCHAR(100) UNIQUE NOT NULL,
        encrypted_pan TEXT NOT NULL,
        last4 VARCHAR(4) NOT NULL,
        brand VARCHAR(20) NOT NULL,
        expiry_month VARCHAR(2) NOT NULL,
        expiry_year VARCHAR(4) NOT NULL,
        merchant_id VARCHAR(36) NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

    CREATE INDEX IF NOT EXISTS idx_tokens_token ON tokens(token);
    CREATE INDEX IF NOT EXISTS idx_tokens_merchant_id ON tokens(merchant_id);
    `
    
    _, err := db.Exec(schema)
    if err != nil {
        return err
    }
    
    log.Println("Database tables created/verified successfully")
    return nil
}
