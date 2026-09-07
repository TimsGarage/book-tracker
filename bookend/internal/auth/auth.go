package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token has expired")
)

// Claims represents the payload of the JWT
type Claims struct {
	UserID   uint   `json:"sub"`
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
}

// jwtHeader defines the standard JWT header for HS256
type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a bcrypt hashed password with its possible plaintext equivalent
func CheckPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// GenerateToken creates a signed standard HMAC-SHA256 (HS256) JWT string
func GenerateToken(userID uint, username string, secret string, ttl time.Duration) (string, error) {
	if secret == "" {
		return "", errors.New("jwt secret cannot be empty")
	}

	header := jwtHeader{
		Alg: "HS256",
		Typ: "JWT",
	}
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to encode jwt header: %w", err)
	}
	encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)

	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Iat:      now.Unix(),
		Exp:      now.Add(ttl).Unix(),
	}
	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to encode jwt claims: %w", err)
	}
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsBytes)

	signingInput := encodedHeader + "." + encodedClaims

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	signature := mac.Sum(nil)
	encodedSignature := base64.RawURLEncoding.EncodeToString(signature)

	return signingInput + "." + encodedSignature, nil
}

// ValidateToken parses and verifies an HS256 JWT string, checking signature and expiration
func ValidateToken(tokenStr string, secret string) (*Claims, error) {
	if secret == "" {
		return nil, errors.New("jwt secret cannot be empty")
	}

	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	encodedHeader := parts[0]
	encodedClaims := parts[1]
	encodedSignature := parts[2]

	// Verify header
	headerBytes, err := base64.RawURLEncoding.DecodeString(encodedHeader)
	if err != nil {
		return nil, ErrInvalidToken
	}
	var header jwtHeader
	if err := json.Unmarshal(headerBytes, &header); err != nil || header.Alg != "HS256" {
		return nil, ErrInvalidToken
	}

	// Verify signature
	signingInput := encodedHeader + "." + encodedClaims
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	expectedSignature := mac.Sum(nil)

	actualSignature, err := base64.RawURLEncoding.DecodeString(encodedSignature)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if !hmac.Equal(actualSignature, expectedSignature) {
		return nil, ErrInvalidToken
	}

	// Decode claims
	claimsBytes, err := base64.RawURLEncoding.DecodeString(encodedClaims)
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if claims.Exp < time.Now().Unix() {
		return nil, ErrTokenExpired
	}

	return &claims, nil
}
