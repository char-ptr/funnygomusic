package entries

import (
	"funnygomusic/bot"
	"funnygomusic/bot/players"
)

type Url struct {
	url string
}

func (l *Url) GetTitle() string {
	return "<URL>"
}
func (l *Url) GetAlbum() string {
	return l.url
}
func (l *Url) GetArtist() string {
	return "unknown"
}
func (l *Url) GetDuration() int {
	return 0
}
func (l *Url) GetPlayer() bot.Player {
	return players.NewYoutubePlayer(l.url)
}

func NewUrl(url string) (idx Url) {
	idx = Url{url: url}
	return
}
