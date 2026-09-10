package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bookend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.CORSMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	t.Run("sets CORS headers on standard request with Origin", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Origin", "http://localhost:1420")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
		if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:1420" {
			t.Errorf("expected Access-Control-Allow-Origin to be http://localhost:1420, got '%s'", got)
		}
	})

	t.Run("handles OPTIONS preflight request", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodOptions, "/test", nil)
		req.Header.Set("Origin", "http://localhost:1420")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected status 204 No Content for OPTIONS, got %d", w.Code)
		}
		if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:1420" {
			t.Errorf("expected Access-Control-Allow-Origin to be http://localhost:1420, got '%s'", got)
		}
		if got := w.Header().Get("Access-Control-Allow-Methods"); got == "" {
			t.Errorf("expected Access-Control-Allow-Methods header, got empty")
		}
	})
}
