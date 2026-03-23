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

func getPlaylistInfo(binPath, url string) ([]PlaylistEntry, error) {
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

func filterPlaylist(entries []PlaylistEntry, songs map[string]bool) []string {
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

func setsEqual(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

func downloadSongs(binPath string, songs []string, path string) error {
	args := []string{
		"-f", "bestaudio/best",
		"-o", fmt.Sprintf("%s/%%(title)s - %%(artist)s.%%(ext)s", path),
		"--extract-audio",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"--embed-thumbnail",
		"--add-metadata",
		"--ignore-errors",
	}
	args = append(args, songs...)

	cmd := exec.Command(binPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func CheckPlaylists(playlists []Playlist, binPath string, playlistStates map[string]map[string]bool) error {
	for _, playlist := range playlists {
		entries, err := getPlaylistInfo(binPath, playlist.URL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting playlist: %v\n", err)
			return err
		}

		currentIDs := make(map[string]bool, len(entries))
		for _, entry := range entries {
			if entry.ID != "" {
				currentIDs[entry.ID] = true
			}
		}

		if prev, ok := playlistStates[playlist.URL]; ok && setsEqual(prev, currentIDs) {
			fmt.Printf("No changes detected in playlist %s\n", playlist.URL)
			continue
		}
		playlistStates[playlist.URL] = currentIDs

		localSongs, err := GetSongsId(playlist.DownloadPath)
		if err != nil {
			return fmt.Errorf("error getting local songs: %w", err)
		}

		missing := filterPlaylist(entries, localSongs)
		if len(missing) == 0 {
			fmt.Println("No new songs to download.")
			return nil
		}

		if err := downloadSongs(binPath, missing, playlist.DownloadPath); err != nil {
			fmt.Fprintf(os.Stderr, "error downloading songs: %v\n", err)
			return err
		}
	}
	return nil
}

func CheckPlaylistsJob(binPath string, playlists []Playlist, playlistStates map[string]map[string]bool) func() error {
	return func() error {
		return CheckPlaylists(playlists, binPath, playlistStates)
	}
}
