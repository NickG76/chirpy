package auth

import (
	"fmt"
	"testing"
	"time"
	"github.com/google/uuid"
)

func TestJWT_RoundTrip(t *testing.T) {
	secret := "supersecret"
	uid := uuid.New()

	token, err := MakeJWT(uid, secret, time.Minute) 
	if err != nil {
		t.Fatalf("Failed to make jwt: %s", err)
	}

	parsedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("failed to validate jwt %v", err)
	}

	if parsedID != uid {
		t.Errorf("expected userID %s, got %s", uid, parsedID)
	}
}
