package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(header http.Header) (string, error) {
	if header == nil {
		return "", fmt.Errorf("header does not exist")
	}
	headerVal := header.Get("Authorization")
	if headerVal == "" {
		return "", fmt.Errorf("header does not exist")
	}
	words := strings.Split(headerVal, " ")
	if len(words) < 2 {
		return "", fmt.Errorf("api key not found")
	}
	apiKey := words[1]
	return apiKey, nil
}
