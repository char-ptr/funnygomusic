package bot

import (
	"context"
	"funnygomusic/databaser"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/voice"
	"gorm.io/gorm"
)

type ComData int

const (
	NewItem ComData = iota
	SongEnded
	PlaySong
	Exit
)

type Botter struct {
	BState     *state.State
	VoiceSes   *voice.Session
	PlayData   *PlayingData
	Ctx        context.Context
	MyId       discord.UserID
	ComChan    chan ComData
	Queue      []QueueEntry
	AllowList  []string
	QueueIndex int
	awaitSong  bool
	SubChan    discord.ChannelID
	Db         *gorm.DB
}

func NewBotter(s *state.State, ctx *context.Context) Botter {
	return Botter{
		BState:    s,
		Ctx:       *ctx,
		ComChan:   make(chan ComData),
		AllowList: []string{},
		awaitSong: true,
		Db:        databaser.NewDatabase(),
	}

}
func (b *Botter) NewPlayData(entry QueueEntry) {
	pd := NewPlayingData(b, entry.Path)
	b.PlayData = pd
}
func (b *Botter) PlayManagerStart() {
	for {
		v := <-b.ComChan
		switch v {
		case NewItem:
			{
				if b.awaitSong {
					b.awaitSong = false
					go b.requestPlaySong()
				}
			}
		case SongEnded:
			{
				b.QueueIndex++
				log.Println("song ended -> idx = ", b.QueueIndex)
				go b.requestPlaySong()
			}
		case Exit:
			{
				return
			}
		case PlaySong:
			{
				if b.QueueIndex >= len(b.Queue) {
					b.awaitSong = true
					b.PlayData = nil
					log.Printf("queue size too small.. [%d:%d]", b.QueueIndex, len(b.Queue))
					b.BState.SendMessage(b.SubChan, "Queue has ended")

					continue
				}
				b.awaitSong = false
				log.Println("request to play song")
				b.NewPlayData(b.Queue[b.QueueIndex])
				b.PlayData.Start()
				go b.PlayData.SendSongInfo()
			}
		}
	}
}
func (b *Botter) ClearQueue() {
	b.Queue = nil
	var newi int
	b.QueueIndex = newi
}
func (b *Botter) requestPlaySong() {
	b.ComChan <- PlaySong
}
func (b *Botter) CurrentPlayingSong() *QueueEntry {
	if b.QueueIndex >= len(b.Queue) {
		return nil
	}
	return &b.Queue[b.QueueIndex]
}
