FROM golang:alpine
LABEL authors="constanna"

WORKDIR /usr/src/funnygomusic

COPY go.mod go.sum ./

COPY . .

RUN go build -v -o /usr/local/bin/funnygomusic

RUN apk add ffmpeg

CMD ["funnygomusic"]