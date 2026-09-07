package handlers

import (
	"bookend/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LookupHandler struct {
	db *gorm.DB
}

func NewLookupHandler(db *gorm.DB) *LookupHandler {
	return &LookupHandler{db: db}
}

// Intermediate structs to decode the Google Books JSON
type GoogleBooksResponse struct {
	Items []struct {
		VolumeInfo struct {
			Title       string   `json:"title"`
			Authors     []string `json:"authors"`
			Description string   `json:"description"`
			PageCount   int      `json:"pageCount"`
			ImageLinks  struct {
				Thumbnail string `json:"thumbnail"`
			} `json:"imageLinks"`
			IndustryIdentifiers []struct {
				Type       string `json:"type"`
				Identifier string `json:"identifier"`
			} `json:"industryIdentifiers"`
		} `json:"volumeInfo"`
	} `json:"items"`
}

// Example isbn: 379200027X
func (h *LookupHandler) IsbnLookup(c *gin.Context) {
	isbn := c.Query("isbn")
	if isbn == "" {
		log.Printf("No isbn provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide a valid isbn"})
		return
	}

	apiKey := os.Getenv("GOOGLE_BOOKS_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query google books"})
		log.Printf("GOOGLE_BOOKS_KEY not set")
		return
	}

	url := "https://www.googleapis.com/books/v1/volumes?q=isbn:" + isbn + "&key=" + apiKey

	reqClient := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	res, getErr := reqClient.Do(req)
	if getErr != nil {
		log.Printf("Error requesting Google Books API: %v", getErr)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to contact Google Books API"})
		return
	}
	defer res.Body.Close()
	var apiRes GoogleBooksResponse
	if err := json.NewDecoder(res.Body).Decode(&apiRes); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API response"})
		return
	}

	if len(apiRes.Items) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	volume := apiRes.Items[0].VolumeInfo

	book := models.LookupBook{
		Title:         volume.Title,
		Author:        strings.Join(volume.Authors, ", "),
		Description:   volume.Description,
		ThumbnailLink: volume.ImageLinks.Thumbnail,
		Pages:         volume.PageCount,
	}

	// Extract the ISBN (Prefer ISBN_13)
	for _, ident := range volume.IndustryIdentifiers {
		if ident.Type == "ISBN_13" {
			book.Isbn = ident.Identifier
			break
		} else if ident.Type == "ISBN_10" && book.Isbn == "" {
			book.Isbn = ident.Identifier
		}
	}

	c.JSON(http.StatusOK, book)
}
