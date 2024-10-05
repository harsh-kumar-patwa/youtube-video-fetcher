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

func (worker *Worker) Start() {
	ticker := time.NewTicker(worker.interval)
	go func() {
		for {
			worker.fetchAndStoreVideos()
			<-ticker.C
		}
	}()
}

func (worker *Worker) fetchAndStoreVideos() {
	log.Println("Fetching new videos...")
	videos,err := worker.client.FetchVideos(worker.query, time.Now().Add(-24*time.Hour))
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
		err := worker.db.InsertVideo(dbVideo)
		if err != nil {
			log.Printf("Error inserting video: %v", err)
		}
	}
}

func parseTime(timeStr string) time.Time {
	t, _ := time.Parse(time.RFC3339, timeStr)
	return t
}