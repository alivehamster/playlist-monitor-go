package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/a-h/templ"
	"github.com/alivehamster/playlist-manager-go/utils"
	"github.com/gofiber/fiber/v3"
)

func stuff() {
	binPath := flag.String("bin", "./yt-dlp_linux", "path to yt-dlp binary")
	playlistURL := flag.String("url", "", "playlist URL")
	outputPath := flag.String("out", "./downloads", "output directory")
	flag.Parse()

	if *playlistURL == "" {
		fmt.Fprintln(os.Stderr, "error: --url is required")
		os.Exit(1)
	}

	entries, err := utils.GetPlaylistInfo(*binPath, *playlistURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting playlist: %v\n", err)
		os.Exit(1)
	}

	// TODO: load known songs from storage
	known := map[string]bool{}

	missing := utils.FilterPlaylist(entries, known)
	if len(missing) == 0 {
		fmt.Println("No new songs to download.")
		return
	}

	if err := utils.DownloadSongs(*binPath, missing, *outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "error downloading songs: %v\n", err)
		os.Exit(1)
	}
}

func render(c fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func main() {

	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		return render(c, body(config.Playlists))
	})

	app.Post("/add-playlist", func(c fiber.Ctx) error {
		url := c.FormValue("url")
		downloadPath := c.FormValue("downloadPath")

		if url == "" || downloadPath == "" {
			return c.Status(fiber.StatusBadRequest).SendString("URL and Download Path are required")
		}

		config.Playlists = append(config.Playlists, utils.Playlist{
			URL:          url,
			DownloadPath: downloadPath,
		})

		if err := utils.SaveConfig(config); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save config")
		}

		return c.Redirect().To("/")
	})

	log.Fatal(app.Listen(":3000"))
}
