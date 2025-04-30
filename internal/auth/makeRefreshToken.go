package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// Generates a random 256-bit hex-encoded string
func MakeRefreshToken() (string, error) {
	randData := make([]byte, 32)

	_, err := rand.Read(randData)
	if err != nil {
		return "", err
	}
	encodedStr := hex.EncodeToString(randData)
	return encodedStr, nil
}
