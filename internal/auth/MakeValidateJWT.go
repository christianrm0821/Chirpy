package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
)

// makes a jwt which is a json web token which allows users to make request only on their behalf
// returns a complete signed string with the specified signing method
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	//creating current time(UTC) and putting it in a jwt time struct
	currentTime := time.Now().UTC()
	currTimeJwt := jwt.NewNumericDate(currentTime)

	//getting the expired time and putting it in the jwt time stuct
	expiredTime := currentTime.Add(expiresIn)
	expiresJwt := jwt.NewNumericDate(expiredTime)

	//creating the claim
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  currTimeJwt,
		ExpiresAt: expiresJwt,
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))
}

// validates JWT using the tokenstring and the token secret.
// returns user id/error
func ValidateJWT(tokenstring, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenstring, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		fmt.Printf("token is invalid or expired: %v\n", err)
		return uuid.Nil, err
	}
	userID, err := token.Claims.GetSubject()
	if err != nil {
		fmt.Printf("issue getting the userID: %v", err)
		return uuid.Nil, err
	}
	userIDType := uuid.MustParse(userID)
	if userIDType == uuid.Nil {
		fmt.Println("issue converting userID from string to uuid")
	}
	return userIDType, nil
}
