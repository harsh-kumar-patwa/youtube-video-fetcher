package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config struct holds the application configuration
type Config struct {
	YouTubeAPIKeys []string `json:"youtube_api_keys"`
	SearchQuery string `json:"search_query"`
	FetchInterval time.Duration `json:"fetch_interval"`
	ServerPort string `json:"server_port"`
}

// LoadConfig reads and parses the configuration file
func LoadConfig(file string) (*Config, error) {
	var config Config

	// Open the configuration file
	configFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	// Create a JSON decoder
	jsonParser := json.NewDecoder(configFile)

	// Decode the JSON into the Config struct
	err = jsonParser.Decode(&config)
	if err != nil {
		return nil, err
	}

	// Convert FetchInterval from seconds to time.Duration
	config.FetchInterval *= time.Second

	return &config, nil
}