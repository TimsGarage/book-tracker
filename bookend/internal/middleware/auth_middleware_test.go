package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bookend/internal/auth"
	"bookend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "middleware-test-secret"

	r := gin.New()
	r.Use(middleware.AuthMiddleware(secret))
	r.GET("/protected", func(c *gin.Context) {
		userID, okID := middleware.GetUserID(c)
		username, okUser := middleware.GetUsername(c)
		if !okID || !okUser {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "context missing"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"userID": userID, "username": username})
	})

	tests := []struct {
		name         string
		authHeader   string
		expectedCode int
	}{
		{
			name:         "missing auth header",
			authHeader:   "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "non bearer prefix",
			authHeader:   "Basic dXNlcjpwYXNz",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "empty token string",
			authHeader:   "Bearer ",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid token signature",
			authHeader:   "Bearer header.payload.badsignature",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Errorf("expected status %d, got %d", tc.expectedCode, w.Code)
			}
		})
	}

	t.Run("valid token", func(t *testing.T) {
		token, err := auth.GenerateToken(99, "alice", secret, time.Hour)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		expectedBody := `{"userID":99,"username":"alice"}`
		if w.Body.String() != expectedBody {
			t.Errorf("expected %s, got %s", expectedBody, w.Body.String())
		}
	})
}
