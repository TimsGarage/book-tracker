package middleware_test

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"bookend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func TestBodyRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.BodyRecoveryMiddleware())

	type LoginPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	r.POST("/login-test", func(c *gin.Context) {
		var p LoginPayload
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"username": p.Username, "password": p.Password})
	})

	t.Run("recovers body from X-Payload header when body is empty", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/login-test", nil)
		b64 := base64.StdEncoding.EncodeToString([]byte(`{"username":"admin","password":"secret"}`))
		req.Header.Set("X-Payload", b64)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		if got := w.Body.String(); got != `{"password":"secret","username":"admin"}` {
			t.Errorf("unexpected body: %s", got)
		}
	})

	t.Run("recovers login credentials from query params when body is empty", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/login-test?username=alice&password=secretpassword", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		if got := w.Body.String(); got != `{"password":"secretpassword","username":"alice"}` {
			t.Errorf("unexpected body: %s", got)
		}
	})
}
