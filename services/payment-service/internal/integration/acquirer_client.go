package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/payments-sandbox/pkg/logging"
)

type AcquirerClient struct {
	baseURL    string
	httpClient *logging.TraceAwareClient
}

func NewAcquirerClient(baseURL string) *AcquirerClient {
	return &AcquirerClient{
		baseURL: baseURL,
		httpClient: logging.NewTraceAwareClient(&http.Client{
			Timeout: 10 * time.Second,
		}),
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

func (c *AcquirerClient) Authorize(ctx context.Context, req AuthRequest) (*AuthResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/acquirer/authorize", c.baseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
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
