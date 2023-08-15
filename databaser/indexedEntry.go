package databaser

import (
	"context"
	"encoding/json"
	"github.com/hako/durafmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func NewIndexEntryFromPath(path string, ctx context.Context) (IndexedSong, error) {
	cmd := exec.CommandContext(ctx, "ffprobe",
		"-v", "error",
		"-show_entries", "format=duration:format_tags=track,title,artist,album,musicbrainz_trackid,musicbrainz_albumid,musicbrainz_artistid,date,genre",
		"-of", "json=c=1",
		path)
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		println("err running ffprobe:", err)
		return IndexedSong{}, err
	}
	var raw RawProbeOutput = RawProbeOutput{}
	err = json.Unmarshal(out, &raw)
	if err != nil {
		println("err parsing ffprobe:", err)
		return IndexedSong{}, err
	}
	return ParseRawOutputToIndex(&raw, path), nil
}

func ParseRawOutputToIndex(raw *RawProbeOutput, path string) IndexedSong {
	plen1m, _ := strconv.ParseFloat(raw.Format.Duration, 64)
	parsedTrack, _ := strconv.Atoi(raw.Format.Tags.Track)
	return IndexedSong{
		Title:  raw.Format.Tags.Title,
		Artist: raw.Format.Tags.ArtID,
		Album:  raw.Format.Tags.AlbID,
		ID:     raw.Format.Tags.TrkID,
		File:   path,
		Length: plen1m,
		Track:  parsedTrack,
	}
}

func (e *IndexedSong) Duration() time.Duration {
	return time.Duration(e.Length) * time.Second
}

func (e *IndexedSong) DurationStr() string {
	dur := durafmt.Parse(e.Duration()).LimitFirstN(2)
	return dur.String()
}
