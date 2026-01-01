package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	secret := "my-super-secret-key"
	userID := uuid.New()
	duration := time.Hour

	// Test: Valid Token
	token, err := MakeJWT(userID, secret, duration)
	if err != nil {
		t.Fatalf("Failed to make JWT: %v", err)
	}

	parsedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate valid JWT: %v", err)
	}
	if parsedID != userID {
		t.Errorf("Expected ID %v, got %v", userID, parsedID)
	}

	// Test: Expired Token
	expiredToken, _ := MakeJWT(userID, secret, -time.Hour)
	_, err = ValidateJWT(expiredToken, secret)
	if err == nil {
		t.Error("Validated an expired token")
	}

	// Test: Wrong Secret
	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Error("Validated token with wrong secret")
	}
}
