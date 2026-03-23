package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/bogem/id3v2/v2"
)

var purlRegex = regexp.MustCompile(`[?&]v=([^&]+)`)

func getID(filepath string) (string, error) {
	tag, err := id3v2.Open(filepath, id3v2.Options{Parse: true})
	if err != nil {
		return "", err
	}
	defer tag.Close()

	for _, frame := range tag.GetFrames(tag.CommonID("TXXX")) {
		txxx, ok := frame.(id3v2.UserDefinedTextFrame)
		if ok && txxx.Description == "purl" {
			if match := purlRegex.FindStringSubmatch(txxx.Value); match != nil {
				return match[1], nil
			}
		}
	}
	return "", nil
}

func GetSongsId(playlistsdir string) (map[string]bool, error) {
	songs := make(map[string]bool)
	if err := os.MkdirAll(playlistsdir, 0755); err != nil {
		return nil, fmt.Errorf("error creating directory: %w", err)
	}
	files, err := os.ReadDir(playlistsdir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		id, err := getID(filepath.Join(playlistsdir, file.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting ID: %v\n", err)
			continue
		}
		if id != "" {
			songs[id] = true
		}
	}
	return songs, nil
}
