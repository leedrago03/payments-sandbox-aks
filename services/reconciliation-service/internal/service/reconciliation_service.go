package service

import (
    "fmt"
    "reconciliation-service/internal/model"
    "reconciliation-service/internal/repository"
    "time"
    
    "github.com/google/uuid"
)

type ReconciliationService struct {
    repo *repository.ReconciliationRepository
}

func NewReconciliationService(repo *repository.ReconciliationRepository) *ReconciliationService {
    return &ReconciliationService{repo: repo}
}

func (s *ReconciliationService) CreateReconciliation(req *model.CreateReportRequest) (*model.ReconciliationSummary, error) {
    reportID := uuid.New().String()
    
    // Get internal payment records
    internalPayments, err := s.repo.GetInternalPayments(req.ReportDate)
    if err != nil {
        return nil, err
    }
    
    // Create report
    report := &model.ReconciliationReport{
        ID:         reportID,
        ReportDate: req.ReportDate,
        Source:     req.Source,
        Status:     model.StatusPending,
        CreatedAt:  time.Now(),
    }
    
    // Reconcile records
    var entries []model.ReconciliationEntry
    matched := 0
    mismatched := 0
    missing := 0
    totalAmount := 0.0
    matchedAmount := 0.0
    discrepancyAmount := 0.0
    
    // Match external with internal
    externalMap := make(map[string]model.ExternalPaymentRecord)
    for _, ext := range req.ExternalData {
        externalMap[ext.ExternalID] = ext
    }
    
    internalMap := make(map[string]model.InternalPaymentRecord)
    for _, internal := range internalPayments {
        internalMap[internal.PaymentID] = internal
    }
    
    // Check matches
    for _, internal := range internalPayments {
        totalAmount += internal.Amount
        
        // Try to find matching external record (simplified matching by amount)
        found := false
        for extID, external := range externalMap {
            if external.Amount == internal.Amount && external.Currency == internal.Currency {
                // Match found
                entry := model.ReconciliationEntry{
                    ID:                uuid.New().String(),
                    ReportID:          reportID,
                    InternalPaymentID: internal.PaymentID,
                    ExternalPaymentID: extID,
                    InternalAmount:    internal.Amount,
                    ExternalAmount:    external.Amount,
                    Currency:          internal.Currency,
                    Status:            model.StatusMatched,
                    CreatedAt:         time.Now(),
                }
                entries = append(entries, entry)
                matched++
                matchedAmount += internal.Amount
                delete(externalMap, extID)
                found = true
                break
            }
        }
        
        if !found {
            // No match - missing in external
            entry := model.ReconciliationEntry{
                ID:                uuid.New().String(),
                ReportID:          reportID,
                InternalPaymentID: internal.PaymentID,
                InternalAmount:    internal.Amount,
                Currency:          internal.Currency,
                Status:            model.StatusMissing,
                DiscrepancyReason: "Payment not found in external records",
                CreatedAt:         time.Now(),
            }
            entries = append(entries, entry)
            missing++
            discrepancyAmount += internal.Amount
        }
    }
    
    // Remaining external records are mismatched
    for extID, external := range externalMap {
        entry := model.ReconciliationEntry{
            ID:                uuid.New().String(),
            ReportID:          reportID,
            ExternalPaymentID: extID,
            ExternalAmount:    external.Amount,
            Currency:          external.Currency,
            Status:            model.StatusMismatched,
            DiscrepancyReason: "External payment not found in internal records",
            CreatedAt:         time.Now(),
        }
        entries = append(entries, entry)
        mismatched++
        discrepancyAmount += external.Amount
    }
    
    // Update report summary
    report.TotalRecords = len(internalPayments) + len(externalMap)
    report.MatchedRecords = matched
    report.MismatchedRecords = mismatched
    report.MissingRecords = missing
    report.TotalAmount = totalAmount
    report.MatchedAmount = matchedAmount
    report.DiscrepancyAmount = discrepancyAmount
    
    if mismatched == 0 && missing == 0 {
        report.Status = model.StatusMatched
    } else {
        report.Status = model.StatusMismatched
    }
    
    completedAt := time.Now()
    report.CompletedAt = &completedAt
    
    // Save report
    if err := s.repo.CreateReport(report); err != nil {
        return nil, err
    }
    
    // Save entries
    for _, entry := range entries {
        if err := s.repo.CreateEntry(&entry); err != nil {
            return nil, fmt.Errorf("failed to create entry: %w", err)
        }
    }
    
    return &model.ReconciliationSummary{
        Report:  *report,
        Entries: entries,
    }, nil
}

func (s *ReconciliationService) GetReport(reportID string) (*model.ReconciliationSummary, error) {
    report, err := s.repo.GetReport(reportID)
    if err != nil {
        return nil, err
    }
    if report == nil {
        return nil, fmt.Errorf("report not found")
    }
    
    entries, err := s.repo.GetEntriesByReportID(reportID)
    if err != nil {
        return nil, err
    }
    
    return &model.ReconciliationSummary{
        Report:  *report,
        Entries: entries,
    }, nil
}

func (s *ReconciliationService) GetAllReports(limit int) ([]model.ReconciliationReport, error) {
    if limit <= 0 {
        limit = 50
    }
    return s.repo.GetAllReports(limit)
}
