package model

import "time"

type Merchant struct {
    ID              string    `json:"id" db:"id"`
    BusinessName    string    `json:"business_name" db:"business_name"`
    Email           string    `json:"email" db:"email"`
    ContactPerson   string    `json:"contact_person" db:"contact_person"`
    Phone           string    `json:"phone" db:"phone"`
    Status          string    `json:"status" db:"status"` // ACTIVE, SUSPENDED, CLOSED
    KYCVerified     bool      `json:"kyc_verified" db:"kyc_verified"`
    SettlementBank  string    `json:"settlement_bank" db:"settlement_bank"`
    AccountNumber   string    `json:"account_number" db:"account_number"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type APIKey struct {
    ID          string    `json:"id" db:"id"`
    MerchantID  string    `json:"merchant_id" db:"merchant_id"`
    KeyHash     string    `json:"key_hash" db:"key_hash"` // bcrypt hash
    KeyPrefix   string    `json:"key_prefix" db:"key_prefix"` // First 8 chars for display
    Name        string    `json:"name" db:"name"` // Key name/description
    Status      string    `json:"status" db:"status"` // ACTIVE, REVOKED
    LastUsedAt  *time.Time `json:"last_used_at" db:"last_used_at"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Settlement struct {
    ID              string    `json:"id" db:"id"`
    MerchantID      string    `json:"merchant_id" db:"merchant_id"`
    Amount          float64   `json:"amount" db:"amount"`
    Currency        string    `json:"currency" db:"currency"`
    Status          string    `json:"status" db:"status"` // PENDING, PROCESSING, COMPLETED, FAILED
    PaymentCount    int       `json:"payment_count" db:"payment_count"`
    SettlementDate  time.Time `json:"settlement_date" db:"settlement_date"`
    CompletedAt     *time.Time `json:"completed_at" db:"completed_at"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// Request/Response DTOs
type CreateMerchantRequest struct {
    BusinessName   string `json:"business_name" validate:"required"`
    Email          string `json:"email" validate:"required,email"`
    ContactPerson  string `json:"contact_person" validate:"required"`
    Phone          string `json:"phone"`
    SettlementBank string `json:"settlement_bank"`
    AccountNumber  string `json:"account_number"`
}

type CreateAPIKeyRequest struct {
    Name string `json:"name" validate:"required"`
}

type CreateAPIKeyResponse struct {
    ID       string `json:"id"`
    APIKey   string `json:"api_key"` // Only shown once!
    KeyPrefix string `json:"key_prefix"`
    Name     string `json:"name"`
}
