package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"bookend/internal/middleware"
	"bookend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BookHandler manages book-related REST API endpoints
type BookHandler struct {
	db *gorm.DB
}

func NewBookHandler(db *gorm.DB) *BookHandler {
	return &BookHandler{db: db}
}

func (h *BookHandler) GetMyBooks(c *gin.Context) {
	var books []models.Book

	query := h.db.Order("id desc")
	if userID, ok := middleware.GetUserID(c); ok && userID != 0 {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  books,
		"count": len(books),
	})
}

func (h *BookHandler) GetWishlistBooks(c *gin.Context) {
	var books []models.Book

	query := h.db.Order("id desc")
	if userID, ok := middleware.GetUserID(c); ok && userID != 0 {
		query = query.Where("user_id = ?", userID)
	}
	query = query.Where("ownership_status = ?", "wishlist")

	if err := query.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  books,
		"count": len(books),
	})
}

// GetBook returns a single book by ID
// GET /api/v1/books/:id
func (h *BookHandler) GetBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book models.Book
	if err := h.db.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// CreateBook adds a new book
// POST /api/v1/books
func (h *BookHandler) CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		if isbn, title := c.Query("isbn"), c.Query("title"); isbn != "" && title != "" {
			pages, _ := strconv.Atoi(c.Query("pages"))
			book = models.Book{
				Isbn:          isbn,
				Title:         title,
				Author:        c.Query("author"),
				Release:       c.Query("release"),
				Publisher:     c.Query("publisher"),
				Description:   c.Query("description"),
				ThumbnailLink: c.Query("thumbnail_link"),
				Pages:         pages,

				OwnershipStatus: c.Query("ownership_status"),
				OwnedSince:      c.Query("owned_since"),
				ReadingStatus:   c.Query("reading_status"),
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if book.Isbn == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ISBN is required"})
		return
	}

	var userID uint
	if id, ok := middleware.GetUserID(c); ok {
		userID = id
	}

	// Check that the ISBN is not already in the user's library
	var existingBook models.Book
	if err := h.db.Where("user_id = ? AND isbn = ?", userID, book.Isbn).First(&existingBook).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Book with this ISBN is already in your library"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking book"})
		return
	}

	// Reset ID and assign authenticated user ID
	book.ID = 0
	book.UserId = userID

	if err := h.db.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Book created successfully",
		"data":    book,
	})
}

// UpdateBook updates an existing book
// PUT /api/v1/books/:id
func (h *BookHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	query := h.db
	if userID, ok := middleware.GetUserID(c); ok {
		query = query.Where("user_id = ?", userID)
	}

	var book models.Book
	if err := query.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
		return
	}

	var input models.Book
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.OwnershipStatus != "" {
		updates["ownership_status"] = input.OwnershipStatus
	}
	if input.OwnedSince != "" {
		updates["owned_since"] = input.OwnedSince
	}
	if input.ReadingStatus != "" {
		updates["reading_status"] = input.ReadingStatus
	}

	if err := h.db.Model(&book).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book updated successfully",
		"data":    book,
	})
}

// DeleteBook removes a book by ID (soft delete via GORM)
// DELETE /api/v1/books/:id
func (h *BookHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	query := h.db
	if userID, ok := middleware.GetUserID(c); ok {
		query = query.Where("user_id = ?", userID)
	}

	var book models.Book
	if err := query.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
		return
	}

	if err := h.db.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
