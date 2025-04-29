package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidJWT(t *testing.T) {
	userID1 := uuid.New()
	tokenString, err := MakeJWT(userID1, "mySecret", time.Hour)
	if err != nil {
		t.Errorf("was not expecting an error but got error: %v", err)
	}

	validatedUserID, err := ValidateJWT(tokenString, "mySecret")
	if err != nil {
		t.Errorf("was not exprecting an error but got error: %v", err)
	}

	if validatedUserID != userID1 {
		t.Errorf("user IDs do not match, was expecting: %v but got %v", userID1, validatedUserID)
	} else {
		t.Logf("userIDs match and they should")
	}

	_, err = ValidateJWT(tokenString, "wrong-secret")
	if err == nil {
		t.Error("Was expecting an error but did not get one")
	}

	expiredToken, err := MakeJWT(userID1, "newSecret", -time.Hour)
	if err != nil {
		t.Errorf("was not expecting an error but got error: %v", err)
	}

	_, err = ValidateJWT(expiredToken, "newSecret")
	if err == nil {
		t.Error("was expecting a time expired error but did not get one")
	}
	t.Log("got a time expired error")
}
