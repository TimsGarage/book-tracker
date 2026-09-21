package handlers

import (
	bookconnector "bookend/internal/book_connector"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LookupHandler struct {
	db *gorm.DB
}

func NewLookupHandler(db *gorm.DB) *LookupHandler {
	return &LookupHandler{db: db}
}

// Example isbn: 9780140328721
func (h *LookupHandler) IsbnLookup(c *gin.Context) {
	isbn := c.Query("isbn")
	if isbn == "" {
		log.Printf("No isbn provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide a valid isbn"})
		return
	}

	var dnbResponse = bookconnector.DnbLookupIsbn(isbn)
	if !dnbResponse.Success {
		log.Printf("%s", "dnbLookup failed: "+dnbResponse.Message)
		c.JSON(dnbResponse.StatusCode, gin.H{"error": dnbResponse.Message})
	} else {
		c.JSON(http.StatusOK, dnbResponse.Book)
	}

}

func (h *LookupHandler) SearchByTerm(c *gin.Context) {
	searchterm := c.Query("searchterm")
	if searchterm == "" {
		log.Printf("No searchterm provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide a valid searchterm"})
		return
	}

	const limit = 20

	var dnbResponse = bookconnector.DnbSearchBooks(searchterm, limit)
	if !dnbResponse.Success {
		log.Printf("%s", "Search on Deutsche Nationalbank failed: "+dnbResponse.Message)
		c.JSON(dnbResponse.StatusCode, gin.H{"error": dnbResponse.Message})
	} else {
		c.JSON(http.StatusOK, dnbResponse.Books)
	}

}
