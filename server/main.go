package main

import (
	"log"
	"time"

	"youtube-video-fetcher/config"
	"youtube-video-fetcher/database"
	"youtube-video-fetcher/models"
	"youtube-video-fetcher/youtube"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewDB("youtube_videos.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	client, err := youtube.NewClient(cfg.YouTubeAPIKey)
	if err != nil {
		log.Fatalf("Failed to create YouTube client: %v", err)
	}

	videos, err := client.FetchVideos("golang", time.Now().Add(-24*time.Hour))
	if err != nil {
		log.Fatalf("Failed to fetch videos: %v", err)
	}

	for _, video := range videos {
		dbVideo := &models.Video{
			ID:           video.Id.VideoId,
			Title:        video.Snippet.Title,
			Description:  video.Snippet.Description,
			PublishedAt:  parseTime(video.Snippet.PublishedAt),
			ThumbnailURL: video.Snippet.Thumbnails.Default.Url,
		}
		err := db.InsertVideo(dbVideo)
		if err != nil {
			log.Printf("Failed to insert video: %v", err)
		} else {
			log.Printf("Inserted video: %s", dbVideo.Title)
		}
	}

	// Fetch and print stored videos
	storedVideos, err := db.GetVideos(10, 0)
	if err != nil {
		log.Fatalf("Failed to get videos: %v", err)
	}

	log.Println("Stored videos:")
	for _, v := range storedVideos {
		log.Printf("- %s (Published: %s)", v.Title, v.PublishedAt)
	}
}

func parseTime(timeStr string) time.Time {
	t, _ := time.Parse(time.RFC3339, timeStr)
	return t
}