package routes

import (
	"bookend/internal/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter initializes middleware, template loading, and application routes
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Instantiate handlers
	healthHandler := handlers.NewHealthHandler()
	bookHandler := handlers.NewBookHandler(db)
	lookupHandler := handlers.NewLookupHandler(db)

	// Health check route
	r.GET("/health", healthHandler.Check)

	// API version 1 route group
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Check)

		books := v1.Group("/books")
		{
			// books.GET("/", bookHandler.GetBooks)
			books.GET("/:id", bookHandler.GetBook)
			// books.POST("", bookHandler.CreateBook)
			// books.PUT("/:id", bookHandler.UpdateBook)
			// books.DELETE("/:id", bookHandler.DeleteBook)
		}

		lookup := v1.Group("/lookup")
		{
			lookup.GET("", lookupHandler.IsbnLookup)
		}
	}

	return r
}
