FROM golang:1.26.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
RUN go build -o app .

FROM alpine:latest

RUN apk add --no-cache ffmpeg deno curl

WORKDIR /app

RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/download/2026.03.17/yt-dlp_musllinux \
    -o ./yt-dlp_linux && chmod +x ./yt-dlp_linux

COPY --from=builder /app/app .

CMD ["./app"]