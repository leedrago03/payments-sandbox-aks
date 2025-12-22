package model

// Card represents a payment card.
type Card struct {
	PAN         string `json:"pan"`
	Brand       string `json:"brand"`
	Last4       string `json:"last4"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
}
