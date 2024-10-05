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
    // Load configuration from config.json file
    cfg, err := config.LoadConfig("config.json")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize the database connection
    db, err := database.NewDB("youtube_videos.db")
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Initialize and set up the API handler
    handler := api.NewHandler(db)
    http.HandleFunc("/videos", handler.GetVideos)

    // Start the API server in a separate goroutine
    go func() {
        log.Printf("Starting API server on %s", cfg.ServerPort)
        log.Fatal(http.ListenAndServe(cfg.ServerPort, nil))
    }()

    // Initialize the YouTube client
    client, err := youtube.NewClient(cfg.YouTubeAPIKeys)
    if err != nil {
        log.Printf("Failed to create YouTube client: %v", err)
        // Note: We're logging the error but not fatally exiting
    }

    // Start the worker if the YouTube client was successfully initialized
    if client != nil {
        w := worker.NewWorker(db, client, cfg.SearchQuery, cfg.FetchInterval)
        go w.Start()
        log.Printf("Worker started with query '%s' and interval %v", cfg.SearchQuery, cfg.FetchInterval)
    } else {
        log.Println("Worker not started due to YouTube client initialization failure")
    }

    // Serve static files from the ./static directory
    http.Handle("/", http.FileServer(http.Dir("./static")))

    // Keep the main goroutine running indefinitely
    select {}
}