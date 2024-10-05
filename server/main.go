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

    // Initialize and start the API server
    handler := api.NewHandler(db)
    http.HandleFunc("/videos", handler.GetVideos)

    // Start the API server in a goroutine
    go func() {
        log.Printf("Starting API server on %s", cfg.ServerPort)
        log.Fatal(http.ListenAndServe(cfg.ServerPort, nil))
    }()

    // Initialize and start the worker
    client, err := youtube.NewClient(cfg.YouTubeAPIKeys)
    if err != nil {
        log.Printf("Failed to create YouTube client: %v", err)
        // Note: We're logging the error but not fatally exiting
    }

    if client != nil {
        w := worker.NewWorker(db, client, cfg.SearchQuery, cfg.FetchInterval)
        go w.Start()
        log.Printf("Worker started with query '%s' and interval %v", cfg.SearchQuery, cfg.FetchInterval)
    } else {
        log.Println("Worker not started due to YouTube client initialization failure")
    }

    // Keep the main goroutine running
    select {}
}