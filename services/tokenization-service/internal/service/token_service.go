package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/payments-sandbox/pkg/crypto"
	"github.com/payments-sandbox/pkg/resilience"
	"github.com/sony/gobreaker"
	"tokenization-service/internal/model"
	"tokenization-service/internal/repository"
)

type TokenService struct {
	repo      *repository.TokenRepository
	encryptor *crypto.Encryptor
	dbBreaker *gobreaker.CircuitBreaker
}

func NewTokenService(repo *repository.TokenRepository, encryptor *crypto.Encryptor) *TokenService {
	return &TokenService{
		repo:      repo,
		encryptor: encryptor,
		dbBreaker: resilience.NewCircuitBreaker(resilience.BreakerConfig{
			Name:      "token-db",
			Threshold: 3,
			Timeout:   5 * time.Second,
		}),
	}
}

func (s *TokenService) Tokenize(req *model.TokenizeRequest) (*model.TokenizeResponse, error) {
	// Validate PAN (basic)
	pan := strings.ReplaceAll(req.PAN, " ", "")
	if len(pan) < 13 || len(pan) > 19 {
		return nil, errors.New("invalid card number")
	}

	// Encrypt PAN
	encryptedPAN, err := s.encryptor.Encrypt(pan)
	if err != nil {
		return nil, err
	}

	// Generate token
	tokenStr, err := generateToken()
	if err != nil {
		return nil, err
    }

	// Detect card brand
	brand := detectBrand(pan)

	// Get last 4 digits
	last4 := pan[len(pan)-4:]

	// Store token
	token := &model.Token{
		Token:        tokenStr,
		EncryptedPAN: encryptedPAN,
		Last4:        last4,
		Brand:        brand,
		ExpiryMonth:  req.ExpiryMonth,
		ExpiryYear:   req.ExpiryYear,
		MerchantID:   req.MerchantID,
	}

	_, err = s.dbBreaker.Execute(func() (interface{}, error) {
		return nil, s.repo.Create(token)
	})
	if err != nil {
		return nil, err
	}

	return &model.TokenizeResponse{
		Token:       tokenStr,
		Last4:       last4,
		Brand:       brand,
		ExpiryMonth: req.ExpiryMonth,
		ExpiryYear:  req.ExpiryYear,
	}, nil
}

func (s *TokenService) GetToken(tokenStr string) (*model.Token, error) {
	result, err := s.dbBreaker.Execute(func() (interface{}, error) {
		return s.repo.FindByToken(tokenStr)
	})
	if err != nil {
		return nil, err
	}

	// FindByToken returns (nil, nil) if not found, so result can be nil
	if result == nil {
		return nil, errors.New("token not found")
	}
	token := result.(*model.Token)
	if token == nil {
		return nil, errors.New("token not found")
	}
	return token, nil
}

func (s *TokenService) DetokenizePAN(tokenStr string) (string, error) {
	result, err := s.dbBreaker.Execute(func() (interface{}, error) {
		return s.repo.FindByToken(tokenStr)
	})
	if err != nil {
		return "", err
	}

	if result == nil {
		return "", errors.New("token not found")
	}
	token := result.(*model.Token)
	if token == nil {
		return "", errors.New("token not found")
	}

	// Decrypt PAN
	pan, err := s.encryptor.Decrypt(token.EncryptedPAN)
	if err != nil {
		return "", err
	}

	return pan, nil
}

func generateToken() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "tok_" + base64.URLEncoding.EncodeToString(b)[:32], nil
}

func detectBrand(pan string) string {
	if strings.HasPrefix(pan, "4") {
        return "VISA"
	} else if strings.HasPrefix(pan, "5") {
		return "MASTERCARD"
	} else if strings.HasPrefix(pan, "3") {
		return "AMEX"
	}
	return "UNKNOWN"
}
