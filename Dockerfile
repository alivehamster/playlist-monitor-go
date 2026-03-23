FROM alpine:latest

RUN apk add --no-cache ffmpeg deno curl

WORKDIR /app

RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/download/2026.03.17/yt-dlp_musllinux \
    -o ./yt-dlp_linux && chmod +x ./yt-dlp_linux

COPY app .

CMD ["./app"]