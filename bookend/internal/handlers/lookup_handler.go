package handlers

import (
	"bookend/internal/models"
	"encoding/json"
	"log"
	"net/http"
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

// Intermediate struct to decode the Open Library Books API JSON
type OpenLibraryBook struct {
	Title         string `json:"title"`
	NumberOfPages int    `json:"number_of_pages"`
	Authors       []struct {
		Name string `json:"name"`
	} `json:"authors"`
	Cover struct {
		Large string `json:"large"`
	} `json:"cover"`
	Identifiers struct {
		Isbn13 []string `json:"isbn_13"`
		Isbn10 []string `json:"isbn_10"`
	} `json:"identifiers"`
	Publishers []struct {
		Name string `json:"name"`
	} `json:"publishers"`
	ReleaseDate string `json:"publish_date"`
}

// Example isbn: 9780140328721
func (h *LookupHandler) IsbnLookup(c *gin.Context) {
	isbn := c.Query("isbn")
	if isbn == "" {
		log.Printf("No isbn provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide a valid isbn"})
		return
	}

	// Use the Books API endpoint with the data format
	bibKey := "ISBN:" + isbn
	url := "https://openlibrary.org/api/books?bibkeys=" + bibKey + "&format=json&jscmd=data"

	reqClient := http.Client{
		Timeout: time.Second * 15,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	res, getErr := reqClient.Do(req)
	if getErr != nil {
		log.Printf("Error requesting Open Library API: %v", getErr)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to contact Open Library API"})
		return
	}
	defer res.Body.Close()

	// The API returns a map with the bibKey as the root property
	var apiRes map[string]OpenLibraryBook
	if err := json.NewDecoder(res.Body).Decode(&apiRes); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API response"})
		return
	}

	// Extract the book object using the dynamic key
	volume, exists := apiRes[bibKey]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Map author names
	var authorNames []string
	for _, author := range volume.Authors {
		authorNames = append(authorNames, author.Name)
	}

	// Map publisher names
	var publisherNames []string
	for _, publisher := range volume.Publishers {
		publisherNames = append(publisherNames, publisher.Name)
	}

	// Map to your target struct
	book := models.LookupBook{
		Title:         volume.Title,
		Author:        strings.Join(authorNames, ", "),
		Release:       volume.ReleaseDate,
		Publisher:     strings.Join(publisherNames, ", "),
		ThumbnailLink: volume.Cover.Large,
		Pages:         volume.NumberOfPages,
	}

	// Extract the ISBN (Prefer ISBN_13)
	if len(volume.Identifiers.Isbn13) > 0 {
		book.Isbn = volume.Identifiers.Isbn13[0]
	} else if len(volume.Identifiers.Isbn10) > 0 {
		book.Isbn = volume.Identifiers.Isbn10[0]
	} else {
		book.Isbn = isbn
	}

	c.JSON(http.StatusOK, book)
}
