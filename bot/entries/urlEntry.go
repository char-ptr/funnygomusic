package entries

import (
	"bufio"
	"bytes"
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
}

func (l *Url) GetTitle() string {
	if l.Title == "" {
		return l.Id
	}
	return l.Title
}
func (l *Url) GetAlbum() string {
	return "unknown"
}
func (l *Url) GetArtist() string {
	return l.Channel
}
func (l *Url) GetDuration() int {
	return int(l.Duration * 1000)
}
func (l *Url) GetPlayer() bot.Player {
	return players.NewYoutubePlayer(l.Url)
}
func (l *Url) GetID() string {
	return l.Id
}

func NewUrl(url string) (idx []Url, err error) {
	artworksDir := filepath.Join(databaser.MakeConfigPath(), "artwork")
	tempCmd := exec.Command("yt-dlp", url,
		"--no-simulate",
		"--print", "\"%(.{id,title,channel,duration,channel,timestamp,webpage_url_domain,webpage_url})j\"",
		"--write-thumbnail", "-o", "thumbnail:%(id)s", "--convert-thumbnails", "webp",
		"--skip-download",
		"-P", artworksDir,
	)
	buf := bytes.Buffer{}
	tempCmd.Stdout = &buf
	tempCmd.Stderr = os.Stderr

	err = tempCmd.Run()
	if err != nil {
		println("err running yt-dlp:", err)
		return
	}
	lines := bufio.NewScanner(&buf)
	for lines.Scan() {
		byt := lines.Bytes()
		actualLine := byt[1 : len(byt)-1]
		log.Println(string(actualLine))
		var ytData YtDLPData
		err = json.Unmarshal(actualLine, &ytData)
		if err != nil {
			log.Println("err unmarshalling yt-dlp data:", err)
			continue
		}
		idx = append(idx, Url{&ytData})
	}
	return
}

type YtDLPData struct {
	Id       string  `json:"id"`
	Title    string  `json:"title"`
	Channel  string  `json:"channel"`
	Duration float64 `json:"duration"`
	Domain   string  `json:"webpage_url_domain"`
	Url      string  `json:"webpage_url"`
}
