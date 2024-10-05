package main

import (
	"log"
	"youtube-video-fetcher/config"
	"youtube-video-fetcher/database"
	"youtube-video-fetcher/youtube"
	"youtube-video-fetcher/worker"
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

	// Create and start the worker
	w := worker.NewWorker(db, client, cfg.SearchQuery, cfg.FetchInterval)
	w.Start()

	log.Printf("Worker started with query '%s' and interval %v. Press Ctrl+C to stop.", cfg.SearchQuery, cfg.FetchInterval)
	// Keep the main goroutine running
	select {}
}