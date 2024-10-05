package models

import (
	"time"
)
type Video struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	PublishedAt     time.Time `json:"published_at"`
	ThumbnailURL    string    `json:"thumbnail_url"`
	CreatedAt       time.Time `json:"created_at"`
}