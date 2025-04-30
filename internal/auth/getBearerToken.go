package auth

import (
	"fmt"
	"net/http"
	"strings"
)

// gets token string from header
func GetBearerToken(header http.Header) (string, error) {
	if header == nil {
		return "", fmt.Errorf("header does not exist")
	}
	headerVal := header.Get("Authorization")
	if headerVal == "" {
		return "", fmt.Errorf("header does not exist")
	}
	words := strings.Split(headerVal, " ")
	if len(words) < 2 {
		return "", fmt.Errorf("token not found")
	}
	tokenStr := words[1]
	return tokenStr, nil
}
