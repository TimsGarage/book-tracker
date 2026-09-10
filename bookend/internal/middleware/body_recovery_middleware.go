package middleware

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BodyRecoveryMiddleware recovers request bodies that were stripped by WebViews (like Android WebView shouldInterceptRequest)
func BodyRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
			// 1. Check if X-Payload header is present (base64 encoded JSON)
			if payload := c.GetHeader("X-Payload"); payload != "" {
				if decoded, err := base64.StdEncoding.DecodeString(payload); err == nil && len(decoded) > 0 {
					c.Request.Body = io.NopCloser(bytes.NewReader(decoded))
					c.Request.ContentLength = int64(len(decoded))
					if c.Request.Header.Get("Content-Type") == "" {
						c.Request.Header.Set("Content-Type", "application/json")
					}
					c.Next()
					return
				}
			}

			// 2. Read existing body
			var bodyBytes []byte
			if c.Request.Body != nil && c.Request.Body != http.NoBody {
				bodyBytes, _ = io.ReadAll(c.Request.Body)
			}

			if len(bodyBytes) > 0 {
				c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			} else {
				// Body was empty or nil! Check if query parameters provide a fallback payload
				if c.Query("username") != "" && c.Query("password") != "" {
					jsonPayload, _ := json.Marshal(gin.H{
						"username": c.Query("username"),
						"password": c.Query("password"),
					})
					c.Request.Body = io.NopCloser(bytes.NewReader(jsonPayload))
					c.Request.ContentLength = int64(len(jsonPayload))
					c.Request.Header.Set("Content-Type", "application/json")
				} else if c.Query("isbn") != "" && c.Query("title") != "" {
					pages, _ := strconv.Atoi(c.Query("pages"))
					jsonPayload, _ := json.Marshal(gin.H{
						"isbn":           c.Query("isbn"),
						"title":          c.Query("title"),
						"author":         c.Query("author"),
						"description":    c.Query("description"),
						"thumbnail_link": c.Query("thumbnail_link"),
						"pages":          pages,
					})
					c.Request.Body = io.NopCloser(bytes.NewReader(jsonPayload))
					c.Request.ContentLength = int64(len(jsonPayload))
					c.Request.Header.Set("Content-Type", "application/json")
				}
			}
		}
		c.Next()
	}
}
