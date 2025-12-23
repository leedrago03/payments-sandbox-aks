package service

import (
    "errors"
    "ledger-service/internal/model"
    "ledger-service/internal/repository"
    "time"
    
    "github.com/google/uuid"
)

type LedgerService struct {
    repo *repository.LedgerRepository
}

func NewLedgerService(repo *repository.LedgerRepository) *LedgerService {
    return &LedgerService{repo: repo}
}

func (s *LedgerService) CreateEntries(req *model.CreateEntriesRequest) (*model.CreateEntriesResponse, error) {
    // Validate double-entry (debits must equal credits)
    if !s.validateDoubleEntry(req.Entries) {
        return nil, errors.New("entries not balanced: debits must equal credits")
    }
    
    entryGroupID := uuid.New().String()
    var entries []model.LedgerEntry
    
    for _, entryReq := range req.Entries {
        entry := model.LedgerEntry{
            ID:           uuid.New().String(),
            EntryGroupID: entryGroupID,
            AccountID:    entryReq.AccountID,
            EntryType:    entryReq.EntryType,
            Amount:       entryReq.Amount,
            Currency:     entryReq.Currency,
            PaymentID:    req.PaymentID,
            Description:  req.Description,
            CreatedAt:    time.Now(),
        }
        entries = append(entries, entry)
    }
    
    if err := s.repo.CreateEntries(entries); err != nil {
        return nil, err
    }
    
    return &model.CreateEntriesResponse{
        EntryGroupID: entryGroupID,
        Entries:      entries,
        IsBalanced:   true,
    }, nil
}

func (s *LedgerService) validateDoubleEntry(entries []model.EntryRequest) bool {
    if len(entries) < 2 {
        return false
    }

    var totalDebits, totalCredits float64
    baseCurrency := entries[0].Currency
    
    for _, entry := range entries {
        // All entries in a transaction must be in the same currency
        if entry.Currency != baseCurrency {
            return false
        }

        if entry.EntryType == model.DEBIT {
            totalDebits += entry.Amount
        } else {
            totalCredits += entry.Amount
        }
    }
    
    // Allow 0.0001 difference for floating point precision
    diff := totalDebits - totalCredits
    if diff < 0 {
        diff = -diff
    }
    return diff <= 0.0001
}

func (s *LedgerService) GetPaymentEntries(paymentID string) ([]model.LedgerEntry, error) {
    return s.repo.GetEntriesByPaymentID(paymentID)
}

func (s *LedgerService) GetAccountBalance(accountID, currency string) (*model.BalanceResponse, error) {
    balance, err := s.repo.GetAccountBalance(accountID, currency)
    if err != nil {
        return nil, err
    }
    
    return &model.BalanceResponse{
        AccountID: balance.AccountID,
        Balance:   balance.Balance,
        Currency:  balance.Currency,
    }, nil
}

func (s *LedgerService) GetAllBalances() ([]model.AccountBalance, error) {
    return s.repo.GetAllBalances()
}
