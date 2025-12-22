package model

// TokenizeRequest is the request to tokenize a card.
type TokenizeRequest struct {
	PAN         string `json:"pan"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	CVV         string `json:"cvv"`
	MerchantID  string `json:"merchant_id"`
}

// TokenizeResponse is the response from tokenizing a card.
type TokenizeResponse struct {
	Token       string `json:"token"`
	Last4       string `json:"last4"`
	Brand       string `json:"brand"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
}
