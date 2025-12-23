package model

import "time"

type EventType string

const (
    // Payment events
    PaymentCreated   EventType = "PAYMENT_CREATED"
    PaymentAuthorized EventType = "PAYMENT_AUTHORIZED"
    PaymentCaptured  EventType = "PAYMENT_CAPTURED"
    PaymentFailed    EventType = "PAYMENT_FAILED"
    PaymentRefunded  EventType = "PAYMENT_REFUNDED"
    
    // Merchant events
    MerchantCreated  EventType = "MERCHANT_CREATED"
    MerchantUpdated  EventType = "MERCHANT_UPDATED"
    MerchantSuspended EventType = "MERCHANT_SUSPENDED"
    
    // API Key events
    APIKeyCreated    EventType = "API_KEY_CREATED"
    APIKeyRevoked    EventType = "API_KEY_REVOKED"
    APIKeyUsed       EventType = "API_KEY_USED"
    
    // Token events
    TokenCreated     EventType = "TOKEN_CREATED"
    TokenUsed        EventType = "TOKEN_USED"
    
    // Ledger events
    LedgerEntryCreated EventType = "LEDGER_ENTRY_CREATED"
    
    // Auth events
    LoginSuccess     EventType = "LOGIN_SUCCESS"
    LoginFailure     EventType = "LOGIN_FAILURE"
    LogoutSuccess    EventType = "LOGOUT_SUCCESS"
)

type AuditLog struct {
    ID          string            `json:"id" db:"id"`
    EventType   EventType         `json:"event_type" db:"event_type"`
    EntityType  string            `json:"entity_type" db:"entity_type"`
    EntityID    string            `json:"entity_id" db:"entity_id"`
    ActorType   string            `json:"actor_type" db:"actor_type"`
    ActorID     string            `json:"actor_id" db:"actor_id"`
    Action      string            `json:"action" db:"action"`
    Details     string            `json:"details" db:"details"`
    IPAddress   string            `json:"ip_address" db:"ip_address"`
    UserAgent   string            `json:"user_agent" db:"user_agent"`
    Success     bool              `json:"success" db:"success"`
    ErrorMsg    string            `json:"error_msg,omitempty" db:"error_msg"`
    Signature   string            `json:"signature" db:"signature"`
    CreatedAt   time.Time         `json:"created_at" db:"created_at"`
}

type CreateAuditLogRequest struct {
    EventType  EventType         `json:"event_type" validate:"required"`
    EntityType string            `json:"entity_type" validate:"required"`
    EntityID   string            `json:"entity_id" validate:"required"`
    ActorType  string            `json:"actor_type"`
    ActorID    string            `json:"actor_id"`
    Action     string            `json:"action" validate:"required"`
    Details    string            `json:"details"`
    IPAddress  string            `json:"ip_address"`
    UserAgent  string            `json:"user_agent"`
    Success    bool              `json:"success"`
    ErrorMsg   string            `json:"error_msg,omitempty"`
}

type AuditLogQueryParams struct {
    EntityType string
    EntityID   string
    ActorID    string
    EventType  EventType
    StartDate  time.Time
    EndDate    time.Time
    Limit      int
}
