package internal

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

func getSongsId(playlistsdir string) (map[string]string, error) {
	songs := make(map[string]string)
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
		path := filepath.Join(playlistsdir, file.Name())
		id, err := getID(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting ID: %v\n", err)
			continue
		}
		if id != "" {
			songs[id] = path
		}
	}
	return songs, nil
}
