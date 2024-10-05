package database

import (
	"database/sql"
	"time"
	"youtube-video-fetcher/models"
	
	_ "github.com/mattn/go-sqlite3"
)

// DB struct wraps the sql.DB to provide custom database methods
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection and initializes the videos table
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Creating videos table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS videos (
			id TEXT PRIMARY KEY,
			title TEXT,
			description TEXT,
			published_at DATETIME,
			thumbnail_url TEXT,
			created_at DATETIME
		)
	`)
	if err!= nil {
		return nil,err
	}

	return &DB{db}, nil
}

// InsertVideo adds a new video to the database or ignores if it already exists
func (db *DB) InsertVideo(video *models.Video) error {
	_, err := db.Exec(`
		INSERT OR IGNORE INTO videos (id, title, description, published_at, thumbnail_url, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, video.ID, video.Title, video.Description, video.PublishedAt, video.ThumbnailURL, time.Now())

	return err
}

// GetVideos retrieves a paginated list of videos from the database
func (db *DB) GetVideos(page, perPage int) ([]*models.Video, error) {
	// Calculate the offset based on the page number and items per page
	offset := (page - 1) * perPage
	
	// Query the database for videos, ordered by published date
	rows, err := db.Query(`
		SELECT id, title, description, published_at, thumbnail_url, created_at
		FROM videos
		ORDER BY published_at DESC
		LIMIT ? OFFSET ?
	`, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and create Video objects
	var videos []*models.Video
	for rows.Next() {
		var video models.Video
		err := rows.Scan(&video.ID, &video.Title, &video.Description, &video.PublishedAt, &video.ThumbnailURL, &video.CreatedAt)
		if err != nil {
			return nil, err
		}
		videos = append(videos, &video)
	}

	return videos, nil
}

// GetTotalVideos returns the total number of videos in the database
func (db *DB) GetTotalVideos() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM videos").Scan(&count)
	return count, err
}