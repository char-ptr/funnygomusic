package databaser

import (
	"encoding/json"
	"github.com/hako/durafmt"
	"gorm.io/gorm"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type IndexEntry struct {
	gorm.Model
	Title     string `gorm:"index"`
	Artist    string `gorm:"index"`
	Album     string `gorm:"index"`
	Path      string
	Length    float64
	CanoTrack int
}

func NewIndexEntryFromPathDnc(path string) IndexEntry {
	r, err := NewIndexEntryFromPath(path)
	if err != nil {
		log.Println("err parsing ffprobe:", err)
	}
	return r
}
func NewIndexEntryFromPath(path string) (IndexEntry, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration:format_tags=track,title,artist,album", "-of", "json=c=1", path)
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		println("err running ffprobe:", err)
		return IndexEntry{}, err
	}
	var raw RawProbeOutput = RawProbeOutput{}
	err = json.Unmarshal(out, &raw)
	if err != nil {
		println("err parsing ffprobe:", err)
		return IndexEntry{}, err
	}
	return ParseRawOutputToIndex(&raw, path), nil
}

func ParseRawOutputToIndex(raw *RawProbeOutput, path string) IndexEntry {
	plen1m, _ := strconv.ParseFloat(raw.Format.Duration, 64)
	parsedTrack, _ := strconv.Atoi(raw.Format.Tags.Track)
	return IndexEntry{
		Title:     raw.Format.Tags.Title,
		Artist:    raw.Format.Tags.Artist,
		Album:     raw.Format.Tags.Album,
		Path:      path,
		Length:    plen1m,
		CanoTrack: parsedTrack,
	}
}

func (e *IndexEntry) Duration() time.Duration {
	return time.Duration(e.Length) * time.Second
}

type RawProbeOutput struct {
	Format struct {
		Duration string `json:"duration"`
		Tags     struct {
			Track  string `json:"track"`
			Title  string `json:"title"`
			Artist string `json:"artist"`
			Album  string `json:"album"`
		} `json:"tags"`
	} `json:"format"`
}

func (e *IndexEntry) DurationStr() string {
	dur := durafmt.Parse(e.Duration()).LimitFirstN(2)
	return dur.String()
}
