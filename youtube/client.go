package youtube

import (
	"context"
	"time"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Client struct {
	services []*youtube.Service
	currentIndex int
}

func NewClient(apiKeys []string) (*Client, error) {
    var services []*youtube.Service
    for _, key := range apiKeys {
        ctx := context.Background()
        service, err := youtube.NewService(ctx, option.WithAPIKey(key))
        if err != nil {
            return nil, err
        }
        services = append(services, service)
    }
    return &Client{services: services}, nil
}

func (client *Client) FetchVideos(query string, publishedAfter time.Time) ([]*youtube.SearchResult, error) {
    var allResults []*youtube.SearchResult
    pageToken := ""
    maxVideos := 100 

    // Continue fetching until we reach the maximum number of videos
    for len(allResults) < maxVideos {
        // Get the current YouTube service from the client's pool
        service := client.services[client.currentIndex]

        // Prepare the search request
        call := service.Search.List([]string{"id", "snippet"}).
            Q(query).
            Type("video").
            Order("date").
            PublishedAfter(publishedAfter.Format(time.RFC3339)).
            MaxResults(50)
			
        // Add page token if we're not on the first page
        if pageToken != "" {
            call = call.PageToken(pageToken)
        }

        // Execute the API call
        response, err := call.Do()
        if err != nil {
            // If there's an error, try the next API key
            client.currentIndex = (client.currentIndex + 1) % len(client.services)
            if client.currentIndex != 0 {
                // Try again with the next key
                continue
            }
            // If we've tried all keys, return the error
            return nil, err
        }

        // Append the results from this page
        allResults = append(allResults, response.Items...)

        // Check if we've reached or exceeded the desired number of videos
        if len(allResults) >= maxVideos {
            allResults = allResults[:maxVideos] // Trim to exact number if we've exceeded
            break
        }

        // Check if there are more pages
        if response.NextPageToken == "" {
            break // No more results available
        }
        pageToken = response.NextPageToken // Set the token for the next page
    }

    return allResults, nil
}