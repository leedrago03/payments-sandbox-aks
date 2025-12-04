package model

import "time"

type EntryType string

const (
    DEBIT  EntryType = "DEBIT"
    CREDIT EntryType = "CREDIT"
)

type LedgerEntry struct {
    ID           string    `json:"id" db:"id"`
    EntryGroupID string    `json:"entry_group_id" db:"entry_group_id"`
    AccountID    string    `json:"account_id" db:"account_id"`
    EntryType    EntryType `json:"entry_type" db:"entry_type"`
    Amount       float64   `json:"amount" db:"amount"`
    Currency     string    `json:"currency" db:"currency"`
    PaymentID    string    `json:"payment_id" db:"payment_id"`
    Description  string    `json:"description" db:"description"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type AccountBalance struct {
    AccountID string  `json:"account_id" db:"account_id"`
    Balance   float64 `json:"balance" db:"balance"`
    Currency  string  `json:"currency" db:"currency"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateEntriesRequest struct {
    PaymentID   string              `json:"payment_id" validate:"required"`
    Description string              `json:"description"`
    Entries     []EntryRequest      `json:"entries" validate:"required,min=2"`
}

type EntryRequest struct {
    AccountID string    `json:"account_id" validate:"required"`
    EntryType EntryType `json:"entry_type" validate:"required"`
    Amount    float64   `json:"amount" validate:"required,gt=0"`
    Currency  string    `json:"currency" validate:"required"`
}

type CreateEntriesResponse struct {
    EntryGroupID string         `json:"entry_group_id"`
    Entries      []LedgerEntry  `json:"entries"`
    IsBalanced   bool           `json:"is_balanced"`
}

type BalanceResponse struct {
    AccountID string  `json:"account_id"`
    Balance   float64 `json:"balance"`
    Currency  string  `json:"currency"`
}
