package model

type AuthRequest struct {
    Token       string  `json:"token" validate:"required"`
    Amount      float64 `json:"amount" validate:"required"`
    Currency    string  `json:"currency" validate:"required"`
    MerchantID  string  `json:"merchant_id" validate:"required"`
}

type AuthResponse struct {
    Status          string  `json:"status"` // APPROVED, DECLINED, TIMEOUT
    TransactionID   string  `json:"transaction_id"`
    AuthCode        string  `json:"auth_code,omitempty"`
    DeclineReason   string  `json:"decline_reason,omitempty"`
    Amount          float64 `json:"amount"`
    Currency        string  `json:"currency"`
}

type CaptureRequest struct {
    TransactionID string  `json:"transaction_id" validate:"required"`
    Amount        float64 `json:"amount" validate:"required"`
}

type CaptureResponse struct {
    Status        string  `json:"status"` // CAPTURED, FAILED
    TransactionID string  `json:"transaction_id"`
    CapturedAmount float64 `json:"captured_amount"`
}

type RefundRequest struct {
    TransactionID string  `json:"transaction_id" validate:"required"`
    Amount        float64 `json:"amount" validate:"required"`
    Reason        string  `json:"reason"`
}

type RefundResponse struct {
    Status        string  `json:"status"` // REFUNDED, FAILED
    TransactionID string  `json:"transaction_id"`
    RefundID      string  `json:"refund_id"`
    RefundAmount  float64 `json:"refund_amount"`
}
