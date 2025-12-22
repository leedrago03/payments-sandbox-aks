package service

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tokenization-service/internal/crypto"
	"tokenization-service/internal/model"
	"tokenization-service/internal/repository"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tokens (
			id TEXT PRIMARY KEY,
			token TEXT NOT NULL,
			encrypted_pan BLOB NOT NULL,
			last4 TEXT NOT NULL,
			brand TEXT NOT NULL,
			expiry_month INTEGER NOT NULL,
			expiry_year INTEGER NOT NULL,
			merchant_id TEXT NOT NULL,
			created_at DATETIME NOT NULL
		);
	`)
	require.NoError(t, err)

	return db
}

func TestTokenService(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// For testing, we use a fixed key instead of KMS
	os.Setenv("AZURE_KEY_VAULT_URI", "") 
	encryptor, err := crypto.NewEncryptor("", "")
	require.NoError(t, err)

	repo := repository.NewTokenRepository(db)
	service := NewTokenService(repo, encryptor)

	req := &model.TokenizeRequest{
		PAN:         "4242424242424242",
		ExpiryMonth: 12,
		ExpiryYear:  2030,
		CVV:         "123",
		MerchantID:  "merch-123",
	}

	// Test Tokenize
	res, err := service.Tokenize(req)
	require.NoError(t, err)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, "4242", res.Last4)
	assert.Equal(t, "VISA", res.Brand)
	assert.Equal(t, 12, res.ExpiryMonth)
	assert.Equal(t, 2030, res.ExpiryYear)

	// Test GetToken
	token, err := service.GetToken(res.Token)
	require.NoError(t, err)
	assert.Equal(t, res.Token, token.Token)
	assert.Equal(t, "4242", token.Last4)

	// Test DetokenizePAN
	pan, err := service.DetokenizePAN(res.Token)
	require.NoError(t, err)
	assert.Equal(t, "4242424242424242", pan)
}
