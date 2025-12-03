package model

import "time"

type Token struct {
    ID            string    `json:"id" db:"id"`
    Token         string    `json:"token" db:"token"`
    EncryptedPAN  string    `json:"-" db:"encrypted_pan"`
    Last4         string    `json:"last4" db:"last4"`
    Brand         string    `json:"brand" db:"brand"`
    ExpiryMonth   string    `json:"expiry_month" db:"expiry_month"`
    ExpiryYear    string    `json:"expiry_year" db:"expiry_year"`
    MerchantID    string    `json:"merchant_id" db:"merchant_id"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type TokenizeRequest struct {
    PAN         string `json:"pan" validate:"required"`
    ExpiryMonth string `json:"expiry_month" validate:"required"`
    ExpiryYear  string `json:"expiry_year" validate:"required"`
    CVV         string `json:"cvv" validate:"required"`
    MerchantID  string `json:"merchant_id" validate:"required"`
}

type TokenizeResponse struct {
    Token       string `json:"token"`
    Last4       string `json:"last4"`
    Brand       string `json:"brand"`
    ExpiryMonth string `json:"expiry_month"`
    ExpiryYear  string `json:"expiry_year"`
}
