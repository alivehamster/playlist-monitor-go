package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/alivehamster/playlist-monitor-go/utils"
	"github.com/gofiber/fiber/v3"
)

func main() {
	binPath := *flag.String("bin", "./yt-dlp_linux", "path to yt-dlp binary")
	flag.Parse()

	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playlistStates := map[string]map[string]bool{}
	utils.StartDaily(ctx, utils.CheckPlaylistsJob(binPath, config, playlistStates))

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		config.RLock()
		playlists := make([]utils.Playlist, len(config.Data.Playlists))
		copy(playlists, config.Data.Playlists)
		config.RUnlock()
		c.Set("Content-Type", "text/html")
		return body(playlists).Render(c.Context(), c.Response().BodyWriter())
	})

	app.Post("/add-playlist", func(c fiber.Ctx) error {
		url := strings.Clone(c.FormValue("url"))
		downloadPath := strings.Clone(c.FormValue("downloadPath"))

		if url == "" || downloadPath == "" {
			return c.Status(fiber.StatusBadRequest).SendString("URL and Download Path are required")
		}

		config.Lock()
		config.Data.Playlists = append(config.Data.Playlists, utils.Playlist{
			URL:          url,
			DownloadPath: downloadPath,
		})
		config.Unlock()

		if err := utils.SaveConfig(config); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save config")
		}

		return c.Redirect().To("/")
	})

	app.Post("/remove-playlist", func(c fiber.Ctx) error {
		url := strings.Clone(c.FormValue("url"))
		if url == "" {
			return c.Status(fiber.StatusBadRequest).SendString("URL is required")
		}

		config.Lock()
		playlists := config.Data.Playlists
		for i, p := range playlists {
			if p.URL == url {
				config.Data.Playlists = append(playlists[:i], playlists[i+1:]...)
				break
			}
		}
		config.Unlock()

		if err := utils.SaveConfig(config); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save config")
		}

		return c.Redirect().To("/")
	})

	log.Fatal(app.Listen(":3000"))
}
