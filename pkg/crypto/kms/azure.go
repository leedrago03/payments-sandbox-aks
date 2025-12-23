package kms

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
)

// KeyVaultClient is a client for interacting with Azure Key Vault.
type KeyVaultClient struct {
	client *azkeys.Client
}

// NewKeyVaultClient creates a new KeyVaultClient.
func NewKeyVaultClient(keyVaultURI string) (*KeyVaultClient, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	client, err := azkeys.NewClient(keyVaultURI, credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create azkeys client: %w", err)
	}

	return &KeyVaultClient{client: client}, nil
}

// GetEncryptionKey retrieves the encryption key from Azure Key Vault.
func (kv *KeyVaultClient) GetEncryptionKey(ctx context.Context, keyName string) ([]byte, error) {
	// In a real scenario, you would fetch the key from Key Vault.
	// For now, we will use a placeholder key for local development.
	if os.Getenv("AZURE_KEY_VAULT_URI") == "" {
		// This is a placeholder for local development and testing.
		// DO NOT use this in production.
		return []byte("a-super-secret-key-for-local-dev"), nil
	}

	key, err := kv.client.GetKey(ctx, keyName, "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get key from Key Vault: %w", err)
	}

	// The key material is in JWK format. We need to extract the raw key value.
	// For an octet key, it's in the 'k' field.
	if key.Key.K == nil {
		return nil, fmt.Errorf("key material is not available")
	}

	return key.Key.K, nil
}
