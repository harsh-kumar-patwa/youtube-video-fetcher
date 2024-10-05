package main

import (
	"log"
	"net/http"
	"youtube-video-fetcher/config"
	"youtube-video-fetcher/database"
	"youtube-video-fetcher/youtube"
	"youtube-video-fetcher/worker"
	"youtube-video-fetcher/api"
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

	client, err := youtube.NewClient(cfg.YouTubeAPIKeys)
	if err != nil {
    log.Fatalf("Failed to create YouTube client: %v", err)
}

	// Create and start the worker
	w := worker.NewWorker(db, client, cfg.SearchQuery, cfg.FetchInterval)
	w.Start()

	// Create API handler
	handler := api.NewHandler(db)

	// Set up HTTP routes
	http.HandleFunc("/videos", handler.GetVideos)

	// Start the HTTP server
	log.Printf("Starting server on %s", cfg.ServerPort)
	log.Printf("Worker started with query '%s' and interval %v", cfg.SearchQuery, cfg.FetchInterval)
	log.Fatal(http.ListenAndServe(cfg.ServerPort, nil))
}