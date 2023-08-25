FROM golang:alpine
LABEL authors="constanna"

RUN apk add --update ffmpeg yt-dlp


WORKDIR /usr/src/funnygomusic

COPY go.mod go.sum ./

COPY . .

RUN go build -v -o /usr/local/bin/funnygomusic

CMD ["funnygomusic"]