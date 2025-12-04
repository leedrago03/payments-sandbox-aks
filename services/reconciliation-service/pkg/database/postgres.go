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
    CREATE TABLE IF NOT EXISTS reconciliation_reports (
        id VARCHAR(36) PRIMARY KEY,
        report_date TIMESTAMP NOT NULL,
        source VARCHAR(100) NOT NULL,
        total_records INTEGER DEFAULT 0,
        matched_records INTEGER DEFAULT 0,
        mismatched_records INTEGER DEFAULT 0,
        missing_records INTEGER DEFAULT 0,
        total_amount DECIMAL(15, 2) DEFAULT 0,
        matched_amount DECIMAL(15, 2) DEFAULT 0,
        discrepancy_amount DECIMAL(15, 2) DEFAULT 0,
        status VARCHAR(20) NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        completed_at TIMESTAMP
    );

    CREATE INDEX IF NOT EXISTS idx_reports_date ON reconciliation_reports(report_date);
    CREATE INDEX IF NOT EXISTS idx_reports_status ON reconciliation_reports(status);

    CREATE TABLE IF NOT EXISTS reconciliation_entries (
        id VARCHAR(36) PRIMARY KEY,
        report_id VARCHAR(36) NOT NULL REFERENCES reconciliation_reports(id) ON DELETE CASCADE,
        internal_payment_id VARCHAR(100),
        external_payment_id VARCHAR(100),
        internal_amount DECIMAL(15, 2) DEFAULT 0,
        external_amount DECIMAL(15, 2) DEFAULT 0,
        currency VARCHAR(3) NOT NULL DEFAULT 'USD',
        status VARCHAR(20) NOT NULL,
        discrepancy_reason TEXT,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

    CREATE INDEX IF NOT EXISTS idx_entries_report_id ON reconciliation_entries(report_id);
    CREATE INDEX IF NOT EXISTS idx_entries_status ON reconciliation_entries(status);
    `
    
    _, err := db.Exec(schema)
    if err != nil {
        return err
    }
    
    log.Println("Database tables created/verified successfully")
    return nil
}
