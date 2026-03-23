package utils

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type Playlist struct {
	URL          string `json:"url"`
	DownloadPath string `json:"download_path"`
}

type ConfigData struct {
	Playlists []Playlist `json:"playlists"`
}

type Config struct {
	sync.RWMutex
	Data ConfigData
}

func LoadConfig() (*Config, error) {
	b, err := os.ReadFile("config/config.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{}, nil
		}
		return nil, err
	}
	var data ConfigData
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	return &Config{Data: data}, nil
}

func SaveConfig(cfg *Config) error {
	cfg.RLock()
	b, err := json.MarshalIndent(cfg.Data, "", "  ")
	cfg.RUnlock()
	if err != nil {
		return err
	}
	if err := os.MkdirAll("config", 0755); err != nil {
		return err
	}
	return os.WriteFile("config/config.json", b, 0644)
}
