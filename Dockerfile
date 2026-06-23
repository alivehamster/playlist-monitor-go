FROM golang:1.26.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go tool templ generate
RUN go build -o ./app ./cmd/

FROM alpine:latest

ARG YTDLP_VERSION=2026.06.09
WORKDIR /app

RUN apk add --no-cache ffmpeg deno curl

RUN curl -fsSL "https://github.com/yt-dlp/yt-dlp/releases/download/${YTDLP_VERSION}/yt-dlp_musllinux" -o ./yt-dlp_linux
RUN chmod +x ./yt-dlp_linux

COPY --from=builder /app/app .

CMD ["./app"]