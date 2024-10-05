package youtube

import (
	"context"
	"time"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"google.golang.org/api/googleapi"
	"fmt"
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
    initialIndex := client.currentIndex
    
    for {
        service := client.services[client.currentIndex]
        call := service.Search.List([]string{"id", "snippet"}).
            Q(query).
            Type("video").
            Order("date").
            PublishedAfter(publishedAfter.Format(time.RFC3339)).
            MaxResults(50)

        if pageToken != "" {
            call = call.PageToken(pageToken)
        }

        response, err := call.Do()
        if err != nil {
            if apiErr, ok := err.(*googleapi.Error); ok && apiErr.Code == 403 {
                // Quota exceeded, try the next API key
                client.currentIndex = (client.currentIndex + 1) % len(client.services)
                if client.currentIndex == initialIndex {
                    // We've tried all keys and they're all exhausted
                    return nil, fmt.Errorf("all API keys exhausted: %v", err)
                }
                // Try again with the next key
                continue
            }
            // If it's not a quota error, return the error
            return nil, err
        }

        allResults = append(allResults, response.Items...)

        // Check if there are more pages
        if response.NextPageToken == "" {
            break
        }
        pageToken = response.NextPageToken

        if len(allResults) >= 100 {  
            break
        }
    }

    return allResults, nil
}