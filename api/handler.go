package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"youtube-video-fetcher/database"
	"youtube-video-fetcher/models"
)

type Handler struct {
	db *database.DB
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
}

func (handler *Handler) GetVideos(writer http.ResponseWriter, request *http.Request) {
	
	page, err := strconv.Atoi(request.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	
	perPage := 20
	videos, err := handler.db.GetVideos(page, perPage)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	totalCount, err := handler.db.GetTotalVideos()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	totalPages := (totalCount + perPage - 1) / perPage

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

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}