package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// Generates a random 256-bit hex-encoded string
func MakeRefreshToken() (string, error) {
	randData := make([]byte, 32)

	rand.Read(randData)
	encodedStr := hex.EncodeToString(randData)
	return encodedStr, nil
}
