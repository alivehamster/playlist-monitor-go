# What is this
Every 24 hours this will scan your youtube music playlists and download any music that has been added

## Webui is a work in progress

# How do I install
```
# pull the image from github container repository

docker pull ghcr.io/alivehamster/playlist-monitor-go:main

# start the container and go to the webui at localhost:3000 to add playlist and download location
docker run -p 3000:3000 -v ./config:/app/config -v /path/on/host:/path/in/container ghcr.io/alivehamster/playlist-monitor-go:main
```
