package main

import (
	"log"

	"bookend/internal/config"
	"bookend/internal/database"
	"bookend/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Set Gin execution mode
	gin.SetMode(cfg.GinMode)

	// 3. Initialize database connection and auto-migration
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 4. Setup Gin router with handlers, templates, and static assets
	router := routes.SetupRouter(db)

	// 5. Start HTTP server
	log.Printf("Server starting on http://localhost:%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
