package config
import (
	"encoding/json"
	"os"
)

type Config struct {
	YouTubeAPIKey string `json:"youtube_api_key"`
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

	return &config, nil
}