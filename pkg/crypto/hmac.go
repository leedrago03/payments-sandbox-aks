package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// GenerateHMAC generates a SHA256 HMAC signature for a message using a key.
func GenerateHMAC(message []byte, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(message)
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyHMAC verifies a SHA256 HMAC signature for a message using a key.
func VerifyHMAC(message []byte, signature string, key []byte) bool {
	expectedSignature := GenerateHMAC(message, key)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
