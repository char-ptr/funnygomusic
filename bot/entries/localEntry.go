package entries

import (
	"context"
	"funnygomusic/bot"
	"funnygomusic/bot/players"
	"funnygomusic/databaser"
	"strconv"
)

type Local struct {
	title    string
	f        string
	artist   string
	album    string
	duration int
}

func (l *Local) GetTitle() string {
	return l.title
}
func (l *Local) GetAlbum() string {
	return l.album
}
func (l *Local) GetArtist() string {
	return l.artist
}
func (l *Local) GetDuration() int {
	return l.duration
}
func (l *Local) GetPlayer() bot.Player {
	return players.NewFilePlayer(l.f)
}
func NewLocalEntry(ffprobe databaser.RawProbeOutput, path string) *Local {
	parseDuration, err := strconv.ParseFloat(ffprobe.Format.Duration, 64)
	if err != nil {
		return nil
	}
	return &Local{
		title:    ffprobe.Format.Tags.Title,
		f:        path,
		artist:   ffprobe.Format.Tags.Artist,
		album:    ffprobe.Format.Tags.Album,
		duration: int(parseDuration),
	}
}
func NewLocalEntryPath(path string) *Local {
	ffprobe, err := databaser.FetchDataForFile(path, context.TODO())
	if err != nil {
		return nil
	}
	return NewLocalEntry(ffprobe, path)
}
