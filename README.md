# YouTube Video Fetcher

This project is a Go-based API that fetches the latest videos from YouTube based on a search query and stores them in a database. It also provides a simple web interface to view and interact with the stored videos.

## Project Structure

- `main.go`: Entry point of the application
- `config/`: Configuration loading and management
- `api/`: API handlers
- `database/`: Database operations
- `models/`: Data models
- `youtube/`: YouTube API client
- `worker/`: Background worker for fetching videos
- `static/`: Static files for the web interface

## Prerequisites

- Go (version 1.16 or later)
- SQLite3
- A Google Developer account and YouTube Data API v3 key

## Setup

1. Clone the repository:
   ```
   git clone https://github.com/harsh-kumar-patwa/youtube-video-fetcher.git
   cd youtube-video-fetcher
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Create a `config.json` file in the project root with the following content:
   ```json
   {
     "youtube_api_keys": ["YOUR_YOUTUBE_API_KEY"],
     "search_query": "cricket",
     "fetch_interval": 10,
     "server_port": ":8080"
   }
   ```
   Replace `YOUR_YOUTUBE_API_KEY` with your actual YouTube Data API v3 key.

4. Create the `static` directory and add the `index.html` file (provided separately) to it.

## Running the Server

1. Start the server:
   ```
   go run main.go
   ```

2. The server will start on `http://localhost:8080` (or the port specified in your config).

3. The worker will automatically start fetching videos based on the configured search query and interval.

## Testing the API

### 1. Fetch Videos

To retrieve stored videos:

```
GET http://localhost:8080/videos?page=1
```

Parameters:
- `page`: Page number (default: 1)

Response:
```json
{
  "videos": [
    {
      "id": "video_id",
      "title": "Video Title",
      "description": "Video Description",
      "published_at": "2024-10-05T14:11:25Z",
      "thumbnail_url": "https://example.com/thumbnail.jpg",
      "created_at": "2024-10-05T20:00:14.072494+05:30"
    },
    ...
  ],
  "total_count": 111,
  "page": 1,
  "per_page": 20,
  "total_pages": 6
}
```

### 2. Web Interface

Open `http://localhost:8080` in your web browser to access the video dashboard. Here you can:
- View all fetched videos
- Search videos by title or description
- Sort videos by publish date or title
- Navigate through pages of results



## Notes

- The server continuously fetches videos in the background based on the configured interval.
- The API uses pagination to manage large numbers of videos efficiently.
- Make sure your YouTube API quota is sufficient for your fetch interval and expected usage.
