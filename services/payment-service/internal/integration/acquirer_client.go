package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AcquirerClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAcquirerClient(baseURL string) *AcquirerClient {
	return &AcquirerClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type AuthRequest struct {
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	CardToken  string  `json:"card_token"`
	MerchantID string  `json:"merchant_id"`
}

type AuthResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	AuthCode      string `json:"auth_code"`
}

func (c *AcquirerClient) Authorize(req AuthRequest) (*AuthResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(fmt.Sprintf("%s/api/acquirer/authorize", c.baseURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("acquirer returned status %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}
