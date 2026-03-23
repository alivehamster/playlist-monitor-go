FROM alpine:latest

RUN apk add --no-cache ffmpeg

WORKDIR /app

COPY app .

CMD ["./app"]