package auth

import (
	"testing"
	"time"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "superSecret123!"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	if hashed == password {
		t.Fatalf("expected hashed password to differ from plaintext")
	}

	if !CheckPassword(hashed, password) {
		t.Fatalf("expected CheckPassword to return true for correct password")
	}

	if CheckPassword(hashed, "wrongPassword") {
		t.Fatalf("expected CheckPassword to return false for incorrect password")
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	secret := "test-secret-key"
	userID := uint(42)
	username := "testuser"
	ttl := 1 * time.Hour

	token, err := GenerateToken(userID, username, secret, ttl)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("unexpected error validating token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected userID %d, got %d", userID, claims.UserID)
	}
	if claims.Username != username {
		t.Errorf("expected username %s, got %s", username, claims.Username)
	}
}

func TestValidateTokenExpired(t *testing.T) {
	secret := "test-secret-key"
	userID := uint(1)
	username := "expireduser"
	// Generate expired token
	token, err := GenerateToken(userID, username, secret, -10*time.Minute)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	_, err = ValidateToken(token, secret)
	if err == nil {
		t.Fatalf("expected error validating expired token, got nil")
	}
	if err != ErrTokenExpired {
		t.Errorf("expected ErrTokenExpired, got %v", err)
	}
}

func TestValidateTokenTampered(t *testing.T) {
	secret := "test-secret-key"
	token, err := GenerateToken(1, "user1", secret, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	// Validate with wrong secret
	_, err = ValidateToken(token, "different-secret")
	if err == nil {
		t.Fatalf("expected error validating with wrong secret, got nil")
	}
	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}

	// Validate tampered token string
	_, err = ValidateToken(token+"tamper", secret)
	if err == nil {
		t.Fatalf("expected error validating tampered token, got nil")
	}
}
