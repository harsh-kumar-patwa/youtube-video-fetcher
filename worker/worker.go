package worker

import (
	"log"
	"time"

	"youtube-video-fetcher/database"
	"youtube-video-fetcher/models"
	"youtube-video-fetcher/youtube"
)

type Worker struct {
	db     *database.DB
	client *youtube.Client
	query  string
	interval time.Duration
}

func NewWorker(db *database.DB, client *youtube.Client, query string, interval time.Duration) *Worker {
	return &Worker{
		db:     db,
		client: client,
		query:  query,
		interval: interval,
	}
}

func (w *Worker) Start() {
	ticker := time.NewTicker(w.interval)
	go func() {
		for {
			w.fetchAndStoreVideos()
			<-ticker.C
		}
	}()
}

func (w *Worker) fetchAndStoreVideos() {
	log.Println("Fetching new videos...")
	videos, err := w.client.FetchVideos(w.query, time.Now().Add(-w.interval))
	if err != nil {
		log.Printf("Error fetching videos: %v", err)
		return
	}

	for _, video := range videos {
		dbVideo := &models.Video{
			ID:           video.Id.VideoId,
			Title:        video.Snippet.Title,
			Description:  video.Snippet.Description,
			PublishedAt:  parseTime(video.Snippet.PublishedAt),
			ThumbnailURL: video.Snippet.Thumbnails.Default.Url,
		}
		err := w.db.InsertVideo(dbVideo)
		if err != nil {
			log.Printf("Error inserting video: %v", err)
		} else {
			log.Printf("Inserted video: %s", dbVideo.Title)
		}
	}
}

func parseTime(timeStr string) time.Time {
	t, _ := time.Parse(time.RFC3339, timeStr)
	return t
}