package kms

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

// KeyVaultClient is a client for interacting with Azure Key Vault.
type KeyVaultClient struct {
	keysClient    *azkeys.Client
	secretsClient *azsecrets.Client
}

// NewKeyVaultClient creates a new KeyVaultClient.
func NewKeyVaultClient(keyVaultURI string) (*KeyVaultClient, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	keysClient, err := azkeys.NewClient(keyVaultURI, credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create azkeys client: %w", err)
	}

	secretsClient, err := azsecrets.NewClient(keyVaultURI, credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create azsecrets client: %w", err)
	}

	return &KeyVaultClient{
		keysClient:    keysClient,
		secretsClient: secretsClient,
	}, nil
}

// GetEncryptionKey retrieves the encryption key from Azure Key Vault.
// It first attempts to fetch it as a Secret (for Symmetric AES keys in Standard Vaults),
// then falls back to fetching it as a Key resource.
func (kv *KeyVaultClient) GetEncryptionKey(ctx context.Context, keyName string) ([]byte, error) {
	// Try fetching as a Secret first (Recommended for Standard tier Symmetric keys)
	secret, err := kv.secretsClient.GetSecret(ctx, keyName, "", nil)
	if err == nil && secret.Value != nil {
		return []byte(*secret.Value), nil
	}

	// Fallback to Key resource (for Premium tier Symmetric keys)
	key, err := kv.keysClient.GetKey(ctx, keyName, "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get key/secret from Key Vault: %w", err)
	}

	if key.Key.K == nil {
		return nil, fmt.Errorf("key material is not available (is this an RSA key?)")
	}

	return key.Key.K, nil
}