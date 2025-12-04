package model

import "time"

type ReconciliationStatus string

const (
    StatusPending    ReconciliationStatus = "PENDING"
    StatusMatched    ReconciliationStatus = "MATCHED"
    StatusMismatched ReconciliationStatus = "MISMATCHED"
    StatusMissing    ReconciliationStatus = "MISSING"
)

type ReconciliationReport struct {
    ID                string               `json:"id" db:"id"`
    ReportDate        time.Time            `json:"report_date" db:"report_date"`
    Source            string               `json:"source" db:"source"`
    TotalRecords      int                  `json:"total_records" db:"total_records"`
    MatchedRecords    int                  `json:"matched_records" db:"matched_records"`
    MismatchedRecords int                  `json:"mismatched_records" db:"mismatched_records"`
    MissingRecords    int                  `json:"missing_records" db:"missing_records"`
    TotalAmount       float64              `json:"total_amount" db:"total_amount"`
    MatchedAmount     float64              `json:"matched_amount" db:"matched_amount"`
    DiscrepancyAmount float64              `json:"discrepancy_amount" db:"discrepancy_amount"`
    Status            ReconciliationStatus `json:"status" db:"status"`
    CreatedAt         time.Time            `json:"created_at" db:"created_at"`
    CompletedAt       *time.Time           `json:"completed_at,omitempty" db:"completed_at"`
}

type ReconciliationEntry struct {
    ID                  string               `json:"id" db:"id"`
    ReportID            string               `json:"report_id" db:"report_id"`
    InternalPaymentID   string               `json:"internal_payment_id" db:"internal_payment_id"`
    ExternalPaymentID   string               `json:"external_payment_id" db:"external_payment_id"`
    InternalAmount      float64              `json:"internal_amount" db:"internal_amount"`
    ExternalAmount      float64              `json:"external_amount" db:"external_amount"`
    Currency            string               `json:"currency" db:"currency"`
    Status              ReconciliationStatus `json:"status" db:"status"`
    DiscrepancyReason   string               `json:"discrepancy_reason,omitempty" db:"discrepancy_reason"`
    CreatedAt           time.Time            `json:"created_at" db:"created_at"`
}

type CreateReportRequest struct {
    Source      string                     `json:"source" validate:"required"`
    ReportDate  time.Time                  `json:"report_date"`
    ExternalData []ExternalPaymentRecord   `json:"external_data" validate:"required"`
}

type ExternalPaymentRecord struct {
    ExternalID string  `json:"external_id"`
    Amount     float64 `json:"amount"`
    Currency   string  `json:"currency"`
    Date       time.Time `json:"date"`
}

type InternalPaymentRecord struct {
    PaymentID string  `json:"payment_id"`
    Amount    float64 `json:"amount"`
    Currency  string  `json:"currency"`
    Date      time.Time `json:"date"`
}

type ReconciliationSummary struct {
    Report  ReconciliationReport  `json:"report"`
    Entries []ReconciliationEntry `json:"entries"`
}
