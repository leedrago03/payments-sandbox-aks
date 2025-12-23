package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TokenizationClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewTokenizationClient(baseURL string) *TokenizationClient {
	return &TokenizationClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type TokenizeRequest struct {
	PAN         string `json:"pan"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	CVV         string `json:"cvv"`
	MerchantID  string `json:"merchant_id"`
}

type TokenizeResponse struct {
	Token       string `json:"token"`
	Last4       string `json:"last4"`
	Brand       string `json:"brand"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
}

func (c *TokenizationClient) Tokenize(req TokenizeRequest) (*TokenizeResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(fmt.Sprintf("%s/v1/tokenize", c.baseURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("tokenization service returned status %d", resp.StatusCode)
	}

	var tokenizeResp TokenizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenizeResp); err != nil {
		return nil, err
	}

	return &tokenizeResp, nil
}
