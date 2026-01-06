package service

import (
	"errors"
	"ledger-service/internal/model"
	"ledger-service/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/payments-sandbox/pkg/resilience"
	"github.com/sony/gobreaker"
)

type LedgerService struct {
	repo      *repository.LedgerRepository
	dbBreaker *gobreaker.CircuitBreaker
}

func NewLedgerService(repo *repository.LedgerRepository) *LedgerService {
	return &LedgerService{
		repo: repo,
		dbBreaker: resilience.NewCircuitBreaker(resilience.BreakerConfig{
			Name:      "ledger-db",
			Threshold: 3,
			Timeout:   5 * time.Second,
		}),
	}
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

	_, err := s.dbBreaker.Execute(func() (interface{}, error) {
		return nil, s.repo.CreateEntries(entries)
	})
	if err != nil {
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
	result, err := s.dbBreaker.Execute(func() (interface{}, error) {
		return s.repo.GetEntriesByPaymentID(paymentID)
	})
	if err != nil {
		return nil, err
	}
	return result.([]model.LedgerEntry), nil
}

func (s *LedgerService) GetAccountBalance(accountID, currency string) (*model.BalanceResponse, error) {
	result, err := s.dbBreaker.Execute(func() (interface{}, error) {
		return s.repo.GetAccountBalance(accountID, currency)
	})
	if err != nil {
		return nil, err
	}

	balance := result.(*model.AccountBalance)
	return &model.BalanceResponse{
		AccountID: balance.AccountID,
		Balance:   balance.Balance,
		Currency:  balance.Currency,
	}, nil
}

func (s *LedgerService) GetAllBalances() ([]model.AccountBalance, error) {
	result, err := s.dbBreaker.Execute(func() (interface{}, error) {
		return s.repo.GetAllBalances()
	})
	if err != nil {
		return nil, err
	}
	return result.([]model.AccountBalance), nil
}
