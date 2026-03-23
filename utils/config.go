package utils

import (
	"encoding/json"
	"errors"
	"os"
)

type Playlist struct {
	URL          string `json:"url"`
	DownloadPath string `json:"download_path"`
}

type Config struct {
	Playlists []Playlist `json:"playlists"`
}

func LoadConfig() (Config, error) {
	b, err := os.ReadFile("config.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, nil
		}
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func SaveConfig(cfg Config) error {
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("config.json", b, 0644)
}
