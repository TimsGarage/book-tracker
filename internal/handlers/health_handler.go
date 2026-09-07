package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles system health status checks
type HealthHandler struct{}

// NewHealthHandler returns a new HealthHandler instance
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check returns 200 OK with server health information
func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "bookend-api",
	})
}
