package model

import "time"

// Token represents a tokenized card.
type Token struct {
	ID           string    `json:"id"`
	Token        string    `json:"token"`
	EncryptedPAN []byte    `json:"-"`
	Brand        string    `json:"brand"`
	Last4        string    `json:"last4"`
	ExpiryMonth  int       `json:"expiry_month"`
	ExpiryYear   int       `json:"expiry_year"`
	MerchantID   string    `json:"merchant_id"`
	CreatedAt    time.Time `json:"-"`
}