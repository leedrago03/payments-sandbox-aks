package repository

import (
    "database/sql"
    "ledger-service/internal/model"
    "time"
)

type LedgerRepository struct {
    db *sql.DB
}

func NewLedgerRepository(db *sql.DB) *LedgerRepository {
    return &LedgerRepository{db: db}
}

func (r *LedgerRepository) CreateEntries(entries []model.LedgerEntry) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    query := `
        INSERT INTO ledger_entries (id, entry_group_id, account_id, entry_type, 
                                    amount, currency, payment_id, description, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
    
    for _, entry := range entries {
        _, err := tx.Exec(query, entry.ID, entry.EntryGroupID, entry.AccountID,
            entry.EntryType, entry.Amount, entry.Currency, entry.PaymentID,
            entry.Description, entry.CreatedAt)
        if err != nil {
            return err
        }
        
        // Update account balance
        if err := r.updateBalance(tx, entry); err != nil {
            return err
        }
    }
    
    return tx.Commit()
}

func (r *LedgerRepository) updateBalance(tx *sql.Tx, entry model.LedgerEntry) error {
    // Get current balance
    var currentBalance float64
    err := tx.QueryRow(`
        SELECT COALESCE(balance, 0) FROM account_balances 
        WHERE account_id = $1 AND currency = $2
    `, entry.AccountID, entry.Currency).Scan(&currentBalance)
    
    if err != nil && err != sql.ErrNoRows {
        return err
    }
    
    // Calculate new balance
    newBalance := currentBalance
    if entry.EntryType == model.DEBIT {
        newBalance -= entry.Amount
    } else {
        newBalance += entry.Amount
    }
    
    // Upsert balance
    _, err = tx.Exec(`
        INSERT INTO account_balances (account_id, balance, currency, updated_at)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (account_id, currency) 
        DO UPDATE SET balance = $2, updated_at = $4
    `, entry.AccountID, newBalance, entry.Currency, time.Now())
    
    return err
}

func (r *LedgerRepository) GetEntriesByPaymentID(paymentID string) ([]model.LedgerEntry, error) {
    query := `SELECT * FROM ledger_entries WHERE payment_id = $1 ORDER BY created_at`
    
    rows, err := r.db.Query(query, paymentID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var entries []model.LedgerEntry
    for rows.Next() {
        var entry model.LedgerEntry
        err := rows.Scan(&entry.ID, &entry.EntryGroupID, &entry.AccountID,
            &entry.EntryType, &entry.Amount, &entry.Currency, &entry.PaymentID,
            &entry.Description, &entry.CreatedAt)
        if err != nil {
            return nil, err
        }
        entries = append(entries, entry)
    }
    
    return entries, nil
}

func (r *LedgerRepository) GetAccountBalance(accountID, currency string) (*model.AccountBalance, error) {
    balance := &model.AccountBalance{}
    err := r.db.QueryRow(`
        SELECT account_id, balance, currency, updated_at 
        FROM account_balances 
        WHERE account_id = $1 AND currency = $2
    `, accountID, currency).Scan(&balance.AccountID, &balance.Balance,
        &balance.Currency, &balance.UpdatedAt)
    
    if err == sql.ErrNoRows {
        return &model.AccountBalance{
            AccountID: accountID,
            Balance:   0,
            Currency:  currency,
            UpdatedAt: time.Now(),
        }, nil
    }
    
    return balance, err
}

func (r *LedgerRepository) GetAllBalances() ([]model.AccountBalance, error) {
    query := `SELECT account_id, balance, currency, updated_at FROM account_balances ORDER BY account_id`
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var balances []model.AccountBalance
    for rows.Next() {
        var balance model.AccountBalance
        err := rows.Scan(&balance.AccountID, &balance.Balance, &balance.Currency, &balance.UpdatedAt)
        if err != nil {
            return nil, err
        }
        balances = append(balances, balance)
    }
    
    return balances, nil
}
