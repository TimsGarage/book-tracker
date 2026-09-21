package bookconnector

import (
	"bookend/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

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

func OpenLibraryLookupIsbn(isbn string) BookLookupResponse {
	// Use the Books API endpoint with the data format
	bibKey := "ISBN:" + isbn
	url := "https://openlibrary.org/api/books?bibkeys=" + bibKey + "&format=json&jscmd=data"

	reqClient := http.Client{
		Timeout: time.Second * 15,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return BookLookupResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create request",
			Book:       models.LookupBook{},
			Success:    false,
		}
	}

	res, getErr := reqClient.Do(req)
	if getErr != nil {
		log.Printf("Error requesting Open Library API: %v", getErr)
		return BookLookupResponse{
			StatusCode: http.StatusBadGateway,
			Message:    "Failed to contact Open Library API",
			Book:       models.LookupBook{},
			Success:    false,
		}
	}
	defer res.Body.Close()

	// The API returns a map with the bibKey as the root property
	var apiRes map[string]OpenLibraryBook
	if err := json.NewDecoder(res.Body).Decode(&apiRes); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return BookLookupResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to parse API response",
			Book:       models.LookupBook{},
			Success:    false,
		}
	}

	// Extract the book object using the dynamic key
	volume, exists := apiRes[bibKey]
	if !exists {
		return BookLookupResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Book not found",
			Book:       models.LookupBook{},
			Success:    false,
		}
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

	return BookLookupResponse{
		StatusCode: http.StatusOK,
		Message:    "",
		Book:       book,
		Success:    true,
	}
}
