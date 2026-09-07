package handlers

import (
	"errors"
	"net/http"
	"strconv"

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

// GetBooks returns all books
// GET /api/v1/books
func (h *BookHandler) GetMyBooks(c *gin.Context) {
	var books []models.Book

	query := h.db.Order("id desc")
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
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
	var input models.LookupBook
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := models.Book{
		Title:         input.Title,
		Author:        input.Author,
		Description:   input.Description,
		ThumbnailLink: input.ThumbnailLink,
		Pages:         input.Pages,
		Owned:         false,
		Read:          false,
	}

	if err := h.db.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Book created successfully",
		"data":    book,
	})
}

// // UpdateBook updates an existing book
// // PUT /api/v1/books/:id
// func (h *BookHandler) UpdateBook(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.ParseUint(idStr, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
// 		return
// 	}

// 	var book models.Book
// 	if err := h.db.First(&book, id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
// 		return
// 	}

// 	var input models.UpdateBookInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	updates := map[string]interface{}{}
// 	if input.Title != "" {
// 		updates["title"] = input.Title
// 	}
// 	if input.Author != "" {
// 		updates["author"] = input.Author
// 	}
// 	if input.Description != "" {
// 		updates["description"] = input.Description
// 	}
// 	if input.Status != "" {
// 		updates["status"] = input.Status
// 	}

// 	if err := h.db.Model(&book).Updates(updates).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Book updated successfully",
// 		"data":    book,
// 	})
// }

// // DeleteBook removes a book by ID (soft delete via GORM)
// // DELETE /api/v1/books/:id
// func (h *BookHandler) DeleteBook(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.ParseUint(idStr, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
// 		return
// 	}

// 	var book models.Book
// 	if err := h.db.First(&book, id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
// 		return
// 	}

// 	if err := h.db.Delete(&book).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
// }
