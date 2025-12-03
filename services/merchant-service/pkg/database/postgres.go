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
    CREATE TABLE IF NOT EXISTS merchants (
        id VARCHAR(36) PRIMARY KEY,
        business_name VARCHAR(255) NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        contact_person VARCHAR(255) NOT NULL,
        phone VARCHAR(50),
        status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
        kyc_verified BOOLEAN DEFAULT FALSE,
        settlement_bank VARCHAR(255),
        account_number VARCHAR(100),
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

    CREATE TABLE IF NOT EXISTS api_keys (
        id VARCHAR(36) PRIMARY KEY,
        merchant_id VARCHAR(36) NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
        key_hash TEXT NOT NULL,
        key_prefix VARCHAR(20) NOT NULL,
        name VARCHAR(100) NOT NULL,
        status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
        last_used_at TIMESTAMP,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

    CREATE INDEX IF NOT EXISTS idx_merchants_email ON merchants(email);
    CREATE INDEX IF NOT EXISTS idx_api_keys_merchant_id ON api_keys(merchant_id);
    CREATE INDEX IF NOT EXISTS idx_api_keys_status ON api_keys(status);

    CREATE TABLE IF NOT EXISTS settlements (
        id VARCHAR(36) PRIMARY KEY,
        merchant_id VARCHAR(36) NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
        amount DECIMAL(15, 2) NOT NULL,
        currency VARCHAR(3) NOT NULL DEFAULT 'USD',
        status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
        payment_count INTEGER DEFAULT 0,
        settlement_date TIMESTAMP NOT NULL,
        completed_at TIMESTAMP,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

    CREATE INDEX IF NOT EXISTS idx_settlements_merchant_id ON settlements(merchant_id);
    CREATE INDEX IF NOT EXISTS idx_settlements_status ON settlements(status);
    `
    
    _, err := db.Exec(schema)
    if err != nil {
        return err
    }
    
    log.Println("Database tables created/verified successfully")
    return nil
}
