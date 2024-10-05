package database

import (
	"database/sql"
	"time"
	"youtube-video-fetcher/models"
	
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

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

func (db *DB) InsertVideo(video *models.Video) error {
	_, err := db.Exec(`
		INSERT OR IGNORE INTO videos (id, title, description, published_at, thumbnail_url, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, video.ID, video.Title, video.Description, video.PublishedAt, video.ThumbnailURL, time.Now())

	return err
}

func (db *DB) GetVideos(page, perPage int) ([]*models.Video, error) {
	offset := (page - 1) * perPage
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

func (db *DB) GetTotalVideos() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM videos").Scan(&count)
	return count, err
}