package crypto

import (
	"context"
	"fmt"
	"tokenization-service/internal/crypto/kms"
)

// Encryptor handles encryption and decryption of data.
type Encryptor struct {
	kmsClient  *kms.KeyVaultClient
	encryptionKey []byte
}

// NewEncryptor creates a new Encryptor.
func NewEncryptor(keyVaultURI, keyName string) (*Encryptor, error) {
	kmsClient, err := kms.NewKeyVaultClient(keyVaultURI)
	if err != nil {
		return nil, fmt.Errorf("failed to create KMS client: %w", err)
	}

	encryptionKey, err := kmsClient.GetEncryptionKey(context.Background(), keyName)
	if err != nil {
		return nil, fmt.Errorf("failed to get encryption key: %w", err)
	}

	return &Encryptor{
		kmsClient:  kmsClient,
		encryptionKey: encryptionKey,
	}, nil
}

// Encrypt encrypts data using AES-GCM.
func (e *Encryptor) Encrypt(plaintext string) ([]byte, error) {
	return Encrypt([]byte(plaintext), e.encryptionKey)
}

// Decrypt decrypts data using AES-GCM.
func (e *Encryptor) Decrypt(ciphertext []byte) (string, error) {
	decrypted, err := Decrypt(ciphertext, e.encryptionKey)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}
