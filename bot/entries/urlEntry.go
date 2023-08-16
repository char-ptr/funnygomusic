package entries

import (
	"encoding/json"
	"funnygomusic/bot"
	"funnygomusic/bot/players"
	"funnygomusic/databaser"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type Url struct {
	*YtDLPData
	url string
}

func (l *Url) GetTitle() string {
	return l.Title
}
func (l *Url) GetAlbum() string {
	return "unknown"
}
func (l *Url) GetArtist() string {
	return l.Channel
}
func (l *Url) GetDuration() int {
	return l.Duration * 1000
}
func (l *Url) GetPlayer() bot.Player {
	return players.NewYoutubePlayer(l.url)
}

func NewUrl(url string) (idx Url) {
	artworksDir := filepath.Join(databaser.MakeConfigPath(), "artworks")
	tempCmd := exec.Command("yt-dlp", url,
		"--no-simulate",
		"--print", "\"%(.{id,title,channel,duration,channel,timestamp,webpage_url_domain})j\"",
		"--write-thumbnail", "-o", "thumbnail:%(id)s", "--convert-thumbnails", "webp",
		"--skip-download",
		"-P", artworksDir,
	)
	tempCmd.Stderr = os.Stderr

	out, err := tempCmd.Output()
	if err != nil {
		println("err running yt-dlp:", err)
		return
	}
	var raw = YtDLPData{}
	err = json.Unmarshal(out[1:len(out)-2], &raw)
	if err != nil {
		log.Println("err parsing yt-dlp:", err, string(out))
		return
	}

	idx = Url{url: url, YtDLPData: &raw}
	log.Println(raw)
	return
}

type YtDLPData struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Channel  string `json:"channel"`
	Duration int    `json:"duration"`
	Domain   string `json:"webpage_url_domain"`
}
