package routes

import (
	"bookend/internal/config"
	"bookend/internal/handlers"
	"bookend/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter initializes middleware, template loading, and application routes
func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Instantiate handlers
	healthHandler := handlers.NewHealthHandler()
	bookHandler := handlers.NewBookHandler(db)
	lookupHandler := handlers.NewLookupHandler(db)
	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)

	// Auth middleware
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)

	// Health check route
	r.GET("/health", healthHandler.Check)

	// API version 1 route group
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Check)

		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", authMiddleware, authHandler.Me)
			auth.POST("/change-password", authMiddleware, authHandler.ChangePassword)
			auth.PUT("/password", authMiddleware, authHandler.ChangePassword)
		}

		// Book routes protected by auth middleware
		books := v1.Group("/books")
		books.Use(authMiddleware)
		{
			books.GET("", bookHandler.GetMyBooks)
			books.GET("/:id", bookHandler.GetBook)
			books.POST("", bookHandler.CreateBook)
		}

		// Also support singular /book with the same auth protection
		book := v1.Group("/book")
		book.Use(authMiddleware)
		{
			book.GET("", bookHandler.GetMyBooks)
			book.GET("/:id", bookHandler.GetBook)
			book.POST("", bookHandler.CreateBook)
		}

		lookup := v1.Group("/lookup")
		{
			lookup.GET("", lookupHandler.IsbnLookup)
		}
	}

	return r
}
