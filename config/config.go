package config
import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	YouTubeAPIKey string `json:"youtube_api_key"`
	SearchQuery   string `json:"search_query"`
	FetchInterval time.Duration `json:"fetch_interval"`
}

func LoadConfig(file string) (*Config, error) {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return nil, err
	}

	// Convert FetchInterval from seconds to time.Duration
	config.FetchInterval *= time.Second
	return &config, nil
}