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
    CREATE TABLE IF NOT EXISTS ledger_entries (
        id VARCHAR(36) PRIMARY KEY,
        entry_group_id VARCHAR(36) NOT NULL,
        account_id VARCHAR(100) NOT NULL,
        entry_type VARCHAR(10) NOT NULL CHECK (entry_type IN ('DEBIT', 'CREDIT')),
        amount DECIMAL(15, 2) NOT NULL,
        currency VARCHAR(3) NOT NULL DEFAULT 'USD',
        payment_id VARCHAR(100) NOT NULL,
        description TEXT,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

    CREATE INDEX IF NOT EXISTS idx_ledger_payment_id ON ledger_entries(payment_id);
    CREATE INDEX IF NOT EXISTS idx_ledger_entry_group ON ledger_entries(entry_group_id);
    CREATE INDEX IF NOT EXISTS idx_ledger_account_id ON ledger_entries(account_id);

    CREATE TABLE IF NOT EXISTS account_balances (
        account_id VARCHAR(100) NOT NULL,
        currency VARCHAR(3) NOT NULL DEFAULT 'USD',
        balance DECIMAL(15, 2) NOT NULL DEFAULT 0,
        updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
        PRIMARY KEY (account_id, currency)
    );

    CREATE INDEX IF NOT EXISTS idx_balances_account_id ON account_balances(account_id);
    `
    
    _, err := db.Exec(schema)
    if err != nil {
        return err
    }
    
    log.Println("Database tables created/verified successfully")
    return nil
}
