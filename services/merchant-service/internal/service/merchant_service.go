package service

import (
    "crypto/rand"
    "encoding/base64"
    "errors"
    "merchant-service/internal/model"
    "merchant-service/internal/repository"
    
    "golang.org/x/crypto/bcrypt"
)

type MerchantService struct {
    repo *repository.MerchantRepository
}

func NewMerchantService(repo *repository.MerchantRepository) *MerchantService {
    return &MerchantService{repo: repo}
}

func (s *MerchantService) CreateMerchant(req *model.CreateMerchantRequest) (*model.Merchant, error) {
    // Check if merchant already exists
    existing, err := s.repo.FindByEmail(req.Email)
    if err != nil {
        return nil, err
    }
    if existing != nil {
        return nil, errors.New("merchant with this email already exists")
    }
    
    merchant := &model.Merchant{
        BusinessName:   req.BusinessName,
        Email:          req.Email,
        ContactPerson:  req.ContactPerson,
        Phone:          req.Phone,
        SettlementBank: req.SettlementBank,
        AccountNumber:  req.AccountNumber,
    }
    
    if err := s.repo.Create(merchant); err != nil {
        return nil, err
    }
    
    return merchant, nil
}

func (s *MerchantService) GetMerchant(id string) (*model.Merchant, error) {
    merchant, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }
    if merchant == nil {
        return nil, errors.New("merchant not found")
    }
    return merchant, nil
}

func (s *MerchantService) UpdateMerchant(id string, req *model.CreateMerchantRequest) (*model.Merchant, error) {
    merchant, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }
    if merchant == nil {
        return nil, errors.New("merchant not found")
    }
    
    merchant.BusinessName = req.BusinessName
    merchant.ContactPerson = req.ContactPerson
    merchant.Phone = req.Phone
    merchant.SettlementBank = req.SettlementBank
    merchant.AccountNumber = req.AccountNumber
    
    if err := s.repo.Update(merchant); err != nil {
        return nil, err
    }
    
    return merchant, nil
}

func (s *MerchantService) CreateAPIKey(merchantID string, req *model.CreateAPIKeyRequest) (*model.CreateAPIKeyResponse, error) {
    // Verify merchant exists
    merchant, err := s.repo.FindByID(merchantID)
    if err != nil {
        return nil, err
    }
    if merchant == nil {
        return nil, errors.New("merchant not found")
    }
    
    // Generate secure API key
    apiKey, err := generateAPIKey()
    if err != nil {
        return nil, err
    }
    
    // Hash the key for storage
    keyHash, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    apiKeyModel := &model.APIKey{
        MerchantID: merchantID,
        KeyHash:    string(keyHash),
        KeyPrefix:  apiKey[:8], // First 8 chars for display
        Name:       req.Name,
    }
    
    if err := s.repo.CreateAPIKey(apiKeyModel); err != nil {
        return nil, err
    }
    
    return &model.CreateAPIKeyResponse{
        ID:        apiKeyModel.ID,
        APIKey:    apiKey, // Only returned once!
        KeyPrefix: apiKeyModel.KeyPrefix,
        Name:      apiKeyModel.Name,
    }, nil
}

func (s *MerchantService) GetAPIKeys(merchantID string) ([]model.APIKey, error) {
    return s.repo.FindAPIKeysByMerchantID(merchantID)
}

// generateAPIKey creates a secure random API key
func generateAPIKey() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return "sk_live_" + base64.URLEncoding.EncodeToString(b)[:32], nil
}
