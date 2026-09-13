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

	// CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Body recovery middleware (recovers empty/stripped bodies from mobile WebViews)
	r.Use(middleware.BodyRecoveryMiddleware())

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
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authMiddleware, authHandler.Register)
			auth.GET("/me", authMiddleware, authHandler.Me)
			auth.POST("/change-password", authMiddleware, authHandler.ChangePassword)
			auth.PUT("/password", authMiddleware, authHandler.ChangePassword)
		}

		// Book routes protected by auth middleware
		books := v1.Group("/books")
		books.Use(authMiddleware)
		{
			books.GET("/", bookHandler.GetMyBooks)
			books.GET("/owned", bookHandler.GetOwnedBooks)
			books.GET("/search", bookHandler.GetReadBooks)
			books.GET("/wishlist", bookHandler.GetWishlistBooks)
			books.GET("/:id", bookHandler.GetBook)
			books.POST("", bookHandler.CreateBook)
			books.DELETE("/:id", bookHandler.DeleteBook)
		}

		lookup := v1.Group("/lookup")
		{
			lookup.GET("", lookupHandler.IsbnLookup)
		}
	}

	return r
}
