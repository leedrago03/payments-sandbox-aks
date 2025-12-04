package repository

import (
    "database/sql"
    "reconciliation-service/internal/model"
    "time"
)

type ReconciliationRepository struct {
    db *sql.DB
}

func NewReconciliationRepository(db *sql.DB) *ReconciliationRepository {
    return &ReconciliationRepository{db: db}
}

func (r *ReconciliationRepository) CreateReport(report *model.ReconciliationReport) error {
    query := `
        INSERT INTO reconciliation_reports (id, report_date, source, total_records,
                                           matched_records, mismatched_records, missing_records,
                                           total_amount, matched_amount, discrepancy_amount,
                                           status, created_at, completed_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    `
    
    _, err := r.db.Exec(query, report.ID, report.ReportDate, report.Source,
        report.TotalRecords, report.MatchedRecords, report.MismatchedRecords,
        report.MissingRecords, report.TotalAmount, report.MatchedAmount,
        report.DiscrepancyAmount, report.Status, report.CreatedAt, report.CompletedAt)
    
    return err
}

func (r *ReconciliationRepository) CreateEntry(entry *model.ReconciliationEntry) error {
    query := `
        INSERT INTO reconciliation_entries (id, report_id, internal_payment_id,
                                           external_payment_id, internal_amount,
                                           external_amount, currency, status,
                                           discrepancy_reason, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `
    
    _, err := r.db.Exec(query, entry.ID, entry.ReportID, entry.InternalPaymentID,
        entry.ExternalPaymentID, entry.InternalAmount, entry.ExternalAmount,
        entry.Currency, entry.Status, entry.DiscrepancyReason, entry.CreatedAt)
    
    return err
}

func (r *ReconciliationRepository) GetReport(reportID string) (*model.ReconciliationReport, error) {
    report := &model.ReconciliationReport{}
    query := `SELECT * FROM reconciliation_reports WHERE id = $1`
    
    err := r.db.QueryRow(query, reportID).Scan(
        &report.ID, &report.ReportDate, &report.Source, &report.TotalRecords,
        &report.MatchedRecords, &report.MismatchedRecords, &report.MissingRecords,
        &report.TotalAmount, &report.MatchedAmount, &report.DiscrepancyAmount,
        &report.Status, &report.CreatedAt, &report.CompletedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    
    return report, err
}

func (r *ReconciliationRepository) GetEntriesByReportID(reportID string) ([]model.ReconciliationEntry, error) {
    query := `SELECT * FROM reconciliation_entries WHERE report_id = $1 ORDER BY created_at`
    
    rows, err := r.db.Query(query, reportID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var entries []model.ReconciliationEntry
    for rows.Next() {
        var entry model.ReconciliationEntry
        err := rows.Scan(&entry.ID, &entry.ReportID, &entry.InternalPaymentID,
            &entry.ExternalPaymentID, &entry.InternalAmount, &entry.ExternalAmount,
            &entry.Currency, &entry.Status, &entry.DiscrepancyReason, &entry.CreatedAt)
        if err != nil {
            return nil, err
        }
        entries = append(entries, entry)
    }
    
    return entries, nil
}

func (r *ReconciliationRepository) GetAllReports(limit int) ([]model.ReconciliationReport, error) {
    query := `SELECT * FROM reconciliation_reports ORDER BY created_at DESC LIMIT $1`
    
    rows, err := r.db.Query(query, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var reports []model.ReconciliationReport
    for rows.Next() {
        var report model.ReconciliationReport
        err := rows.Scan(&report.ID, &report.ReportDate, &report.Source,
            &report.TotalRecords, &report.MatchedRecords, &report.MismatchedRecords,
            &report.MissingRecords, &report.TotalAmount, &report.MatchedAmount,
            &report.DiscrepancyAmount, &report.Status, &report.CreatedAt, &report.CompletedAt)
        if err != nil {
            return nil, err
        }
        reports = append(reports, report)
    }
    
    return reports, nil
}

func (r *ReconciliationRepository) GetInternalPayments(date time.Time) ([]model.InternalPaymentRecord, error) {
    // Simulated - in production, this would query the payment service database
    // For now, return mock data
    return []model.InternalPaymentRecord{
        {PaymentID: "pay_001", Amount: 100.00, Currency: "USD", Date: date},
        {PaymentID: "pay_002", Amount: 250.50, Currency: "USD", Date: date},
        {PaymentID: "pay_003", Amount: 75.25, Currency: "USD", Date: date},
    }, nil
}
