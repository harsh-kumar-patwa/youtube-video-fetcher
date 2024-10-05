package youtube

import (
	"log"
	"context"
	"time"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Client struct {
	service *youtube.Service
}

func NewClient(apiKey string) (*Client, error) {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &Client{service: service}, nil
}

func (client *Client) FetchVideos(query string, publishedAfter time.Time) ([]*youtube.SearchResult, error) {
	log.Printf("Fetching videos for query: %s, published after: %s", query, publishedAfter.Format(time.RFC3339))
	call := client.service.Search.List([]string{"id", "snippet"}).
		Q(query).
		Type("video").
		Order("date").
		PublishedAfter(publishedAfter.Format(time.RFC3339)).
		MaxResults(50)

	response, err := call.Do()
	if err != nil {
		log.Printf("Error fetching videos: %v", err)
		return nil, err
	}
	return response.Items, nil
}