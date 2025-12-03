package repository

import (
    "database/sql"
    "merchant-service/internal/model"
    "time"
    
    "github.com/google/uuid"
)

type MerchantRepository struct {
    db *sql.DB
}

func NewMerchantRepository(db *sql.DB) *MerchantRepository {
    return &MerchantRepository{db: db}
}

func (r *MerchantRepository) Create(merchant *model.Merchant) error {
    merchant.ID = uuid.New().String()
    merchant.Status = "ACTIVE"
    merchant.KYCVerified = false
    merchant.CreatedAt = time.Now()
    merchant.UpdatedAt = time.Now()
    
    query := `
        INSERT INTO merchants (id, business_name, email, contact_person, phone, 
                             status, kyc_verified, settlement_bank, account_number, 
                             created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `
    
    _, err := r.db.Exec(query, merchant.ID, merchant.BusinessName, merchant.Email,
        merchant.ContactPerson, merchant.Phone, merchant.Status, merchant.KYCVerified,
        merchant.SettlementBank, merchant.AccountNumber, merchant.CreatedAt, merchant.UpdatedAt)
    
    return err
}

func (r *MerchantRepository) FindByID(id string) (*model.Merchant, error) {
    merchant := &model.Merchant{}
    query := `SELECT * FROM merchants WHERE id = $1`
    
    err := r.db.QueryRow(query, id).Scan(
        &merchant.ID, &merchant.BusinessName, &merchant.Email, &merchant.ContactPerson,
        &merchant.Phone, &merchant.Status, &merchant.KYCVerified, &merchant.SettlementBank,
        &merchant.AccountNumber, &merchant.CreatedAt, &merchant.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return merchant, err
}

func (r *MerchantRepository) FindByEmail(email string) (*model.Merchant, error) {
    merchant := &model.Merchant{}
    query := `SELECT * FROM merchants WHERE email = $1`
    
    err := r.db.QueryRow(query, email).Scan(
        &merchant.ID, &merchant.BusinessName, &merchant.Email, &merchant.ContactPerson,
        &merchant.Phone, &merchant.Status, &merchant.KYCVerified, &merchant.SettlementBank,
        &merchant.AccountNumber, &merchant.CreatedAt, &merchant.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return merchant, err
}

func (r *MerchantRepository) Update(merchant *model.Merchant) error {
    merchant.UpdatedAt = time.Now()
    
    query := `
        UPDATE merchants 
        SET business_name = $1, contact_person = $2, phone = $3, 
            settlement_bank = $4, account_number = $5, updated_at = $6
        WHERE id = $7
    `
    
    _, err := r.db.Exec(query, merchant.BusinessName, merchant.ContactPerson,
        merchant.Phone, merchant.SettlementBank, merchant.AccountNumber,
        merchant.UpdatedAt, merchant.ID)
    
    return err
}

func (r *MerchantRepository) CreateAPIKey(apiKey *model.APIKey) error {
    apiKey.ID = uuid.New().String()
    apiKey.Status = "ACTIVE"
    apiKey.CreatedAt = time.Now()
    
    query := `
        INSERT INTO api_keys (id, merchant_id, key_hash, key_prefix, name, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
    
    _, err := r.db.Exec(query, apiKey.ID, apiKey.MerchantID, apiKey.KeyHash,
        apiKey.KeyPrefix, apiKey.Name, apiKey.Status, apiKey.CreatedAt)
    
    return err
}

func (r *MerchantRepository) FindAPIKeysByMerchantID(merchantID string) ([]model.APIKey, error) {
    query := `SELECT id, merchant_id, key_hash, key_prefix, name, status, last_used_at, created_at 
              FROM api_keys WHERE merchant_id = $1 AND status = 'ACTIVE'`
    
    rows, err := r.db.Query(query, merchantID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var keys []model.APIKey
    for rows.Next() {
        var key model.APIKey
        err := rows.Scan(&key.ID, &key.MerchantID, &key.KeyHash, &key.KeyPrefix,
            &key.Name, &key.Status, &key.LastUsedAt, &key.CreatedAt)
        if err != nil {
            return nil, err
        }
        keys = append(keys, key)
    }
    
    return keys, nil
}
