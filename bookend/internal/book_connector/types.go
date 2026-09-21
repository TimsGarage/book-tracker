package bookconnector

import "bookend/internal/models"

type BookLookupResponse struct {
	StatusCode int
	Message    string
	Success    bool
	Book       models.LookupBook
}

type BookSearchResponse struct {
	StatusCode int                 `json:"statusCode"`
	Message    string              `json:"message"`
	Books      []models.LookupBook `json:"books"`
	TotalFound int                 `json:"totalFound"`
	Success    bool                `json:"success"`
}
