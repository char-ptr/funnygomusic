package bot

import (
	"context"
	"io"
	"reflect"
	"time"
)

type QueueEntry interface {
	GetTitle() string
	GetAlbum() string
	GetArtist() string
	GetDuration() int

	GetPlayer() Player
}
type Player interface {
	Play(ctx context.Context) (io.Reader, error)
	Pause() error
	Resume() (io.Reader, error)
	Stop() error
	Seek(ms int) (io.Reader, error)
	Position() int
	PositionDuration() time.Duration
	State() PlayingState
	WaitTillEnd() error
}
type PlayingState int

const (
	PSPlaying PlayingState = iota
	PSPaused
	PSComplete
	PSRestart
	PSNotPlaying
)

func GetTypedEntry(entry *QueueEntry) *TypedEntry {
	return &TypedEntry{
		QueueEntry: *entry,
		Ty:         reflect.TypeOf(*entry).String(),
	}
}

type TypedEntry struct {
	QueueEntry
	Ty string `json:"type"`
}
