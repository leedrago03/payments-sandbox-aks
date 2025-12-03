package repository

import (
    "database/sql"
    "time"
    "tokenization-service/internal/model"
    
    "github.com/google/uuid"
)

type TokenRepository struct {
    db *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
    return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(token *model.Token) error {
    token.ID = uuid.New().String()
    token.CreatedAt = time.Now()
    
    query := `
        INSERT INTO tokens (id, token, encrypted_pan, last4, brand, 
                          expiry_month, expiry_year, merchant_id, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
    
    _, err := r.db.Exec(query, token.ID, token.Token, token.EncryptedPAN,
        token.Last4, token.Brand, token.ExpiryMonth, token.ExpiryYear,
        token.MerchantID, token.CreatedAt)
    
    return err
}

func (r *TokenRepository) FindByToken(tokenStr string) (*model.Token, error) {
    token := &model.Token{}
    query := `SELECT * FROM tokens WHERE token = $1`
    
    err := r.db.QueryRow(query, tokenStr).Scan(
        &token.ID, &token.Token, &token.EncryptedPAN, &token.Last4,
        &token.Brand, &token.ExpiryMonth, &token.ExpiryYear,
        &token.MerchantID, &token.CreatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return token, err
}
