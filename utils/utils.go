package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type PlaylistEntry struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type PlaylistInfo struct {
	Entries []PlaylistEntry `json:"entries"`
}

func GetPlaylistInfo(binPath, url string) ([]PlaylistEntry, error) {
	cmd := exec.Command(binPath, "--flat-playlist", "-J", url)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp error: %w", err)
	}

	var info PlaylistInfo
	if err := json.Unmarshal(out, &info); err != nil {
		return nil, fmt.Errorf("json parse error: %w", err)
	}

	return info.Entries, nil
}

func FilterPlaylist(entries []PlaylistEntry, songs map[string]bool) []string {
	var missing []string

	for _, entry := range entries {
		if songs[entry.ID] {
			continue
		}
		fmt.Printf("Missing: %s\n", entry.Title)
		missing = append(missing, entry.URL)
	}

	return missing
}

func DownloadSongs(binPath string, songs []string, path string) error {
	args := []string{
		"-f", "bestaudio/best",
		"-o", fmt.Sprintf("%s/%%(title)s - %%(artist)s.%%(ext)s", path),
		"--extract-audio",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"--embed-thumbnail",
		"--add-metadata",
		"--write-thumbnail",
		"--ignore-errors",
	}
	args = append(args, songs...)

	cmd := exec.Command(binPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
