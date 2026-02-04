package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key-12345"
	expiresIn := time.Hour

	// Create a JWT
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	if token == "" {
		t.Fatal("MakeJWT returned empty string")
	}

	// Validate the JWT
	parsedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if parsedUserID != userID {
		t.Errorf("ValidateJWT returned wrong user ID: got %v, want %v", parsedUserID, userID)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key-12345"
	expiresIn := -time.Hour // Negative duration means already expired

	// Create an expired JWT
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	// Validate should fail for expired token
	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Error("ValidateJWT should fail for expired token")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "correct-secret-key"
	wrongSecret := "wrong-secret-key"
	expiresIn := time.Hour

	// Create a JWT with correct secret
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	// Validate with wrong secret should fail
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("ValidateJWT should fail when using wrong secret")
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	tokenSecret := "test-secret-key-12345"
	invalidToken := "not-a-valid-jwt"

	_, err := ValidateJWT(invalidToken, tokenSecret)
	if err == nil {
		t.Error("ValidateJWT should fail for invalid token string")
	}
}

func TestMakeJWT_DifferentUsers(t *testing.T) {
	userID1 := uuid.New()
	userID2 := uuid.New()
	tokenSecret := "test-secret-key-12345"
	expiresIn := time.Hour

	token1, err := MakeJWT(userID1, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	token2, err := MakeJWT(userID2, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	// Tokens should be different for different users
	if token1 == token2 {
		t.Error("JWTs for different users should be different")
	}

	// Validate each token returns the correct user ID
	parsedID1, err := ValidateJWT(token1, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}
	if parsedID1 != userID1 {
		t.Errorf("ValidateJWT returned wrong user ID: got %v, want %v", parsedID1, userID1)
	}

	parsedID2, err := ValidateJWT(token2, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}
	if parsedID2 != userID2 {
		t.Errorf("ValidateJWT returned wrong user ID: got %v, want %v", parsedID2, userID2)
	}
}
