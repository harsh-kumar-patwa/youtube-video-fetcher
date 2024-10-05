package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"youtube-video-fetcher/database"
	"youtube-video-fetcher/models"
)

// Handler struct holds a database connection
type Handler struct {
	db *database.DB
}

// NewHandler creates a new Handler instance with a database connection
func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
}

// GetVideos handles the HTTP request for fetching videos with pagination
func (handler *Handler) GetVideos(writer http.ResponseWriter, request *http.Request) {
	
	// Parse the page number from the query parameters, default to 1 if invalid
	page, err := strconv.Atoi(request.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	
	// Set the number of items per page
	perPage := 20

	// Fetch videos from the database for the current page
	videos, err := handler.db.GetVideos(page, perPage)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the total count of videos in the database
	totalCount, err := handler.db.GetTotalVideos()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate the total number of pages
	totalPages := (totalCount + perPage - 1) / perPage

	// Prepare the response structure
	response := struct {
		Videos []*models.Video `json:"videos"`
		TotalCount int `json:"total_count"`
		Page int `json:"page"`
		PerPage int `json:"per_page"`
		TotalPages int `json:"total_pages"`
	}{
		Videos: videos,
		TotalCount: totalCount,
		Page: page,
		PerPage: perPage,
		TotalPages: totalPages,
	}

	// Set the Content-Type header and encode the response as JSON
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}