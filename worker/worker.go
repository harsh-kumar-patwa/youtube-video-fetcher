package worker

import (
	"log"
	"time"

	"youtube-video-fetcher/database"
	"youtube-video-fetcher/models"
	"youtube-video-fetcher/youtube"
)

// Worker struct holds the necessary components for fetching and storing YouTube videos
type Worker struct {
	db       *database.DB
	client   *youtube.Client
	query    string
	interval time.Duration
}

// NewWorker creates a new Worker instance
func NewWorker(db *database.DB, client *youtube.Client, query string, interval time.Duration) *Worker {
	return &Worker{
		db:       db,
		client:   client,
		query:    query,
		interval: interval,
	}
}

// Start begins the worker's main loop to fetch and store videos periodically
func (worker *Worker) Start() {
	ticker := time.NewTicker(worker.interval)
	go func() {
		for {
			worker.fetchAndStoreVideos()
			<-ticker.C // Wait for the next tick before running again
		}
	}()
}

// fetchAndStoreVideos fetches new videos from YouTube and stores them in the database
func (worker *Worker) fetchAndStoreVideos() {
	log.Println("Fetching new videos...")
	// Fetch videos published in the last 24 hours
	videos, err := worker.client.FetchVideos(worker.query, time.Now().Add(-24*time.Hour))
	if err != nil {
		log.Printf("Error fetching videos: %v", err)
		return
	}

	// Iterate through fetched videos and store them in the database
	for _, video := range videos {
		dbVideo := &models.Video{
			ID:           video.Id.VideoId,
			Title:        video.Snippet.Title,
			Description:  video.Snippet.Description,
			PublishedAt:  parseTime(video.Snippet.PublishedAt),
			ThumbnailURL: video.Snippet.Thumbnails.Default.Url,
		}
		err := worker.db.InsertVideo(dbVideo)
		if err != nil {
			log.Printf("Error inserting video: %v", err)
		}
	}
}

// parseTime converts a string time to a time.Time object
func parseTime(timeStr string) time.Time {
	t, _ := time.Parse(time.RFC3339, timeStr)
	return t
}
