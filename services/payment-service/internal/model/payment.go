package model

import "time"

type PaymentStatus string

const (
	StatusAuthorized PaymentStatus = "AUTHORIZED"
	StatusCaptured   PaymentStatus = "CAPTURED"
	StatusFailed     PaymentStatus = "FAILED"
)

type Payment struct {
	ID             string        `json:"id"`
	MerchantID     string        `json:"merchant_id"`
	Amount         float64       `json:"amount"`
	Currency       string        `json:"currency"`
	Token          string        `json:"token"`
	Status         PaymentStatus `json:"status"`
	IdempotencyKey string        `json:"idempotency_key,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type CreatePaymentRequest struct {
	MerchantID     string  `json:"merchant_id"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	Token          string  `json:"token"`
	IdempotencyKey string  `json:"idempotency_key"`
}

type CapturePaymentRequest struct {
	Amount float64 `json:"amount,omitempty"` // Optional partial capture
}
