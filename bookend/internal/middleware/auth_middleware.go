package middleware

import (
	"net/http"
	"strings"

	"bookend/internal/auth"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserIDKey   = "userID"
	ContextUsernameKey = "username"
)

// AuthMiddleware creates a Gin middleware that verifies Bearer JWT tokens
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token string is empty"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(tokenString, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token: " + err.Error()})
			c.Abort()
			return
		}

		// Store user identity in context
		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUsernameKey, claims.Username)

		c.Next()
	}
}

// GetUserID retrieves the authenticated user's ID from the Gin context
func GetUserID(c *gin.Context) (uint, bool) {
	val, exists := c.Get(ContextUserIDKey)
	if !exists {
		return 0, false
	}
	id, ok := val.(uint)
	return id, ok
}

// GetUsername retrieves the authenticated username from the Gin context
func GetUsername(c *gin.Context) (string, bool) {
	val, exists := c.Get(ContextUsernameKey)
	if !exists {
		return "", false
	}
	username, ok := val.(string)
	return username, ok
}
