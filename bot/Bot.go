package bot

import (
	"context"
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/voice"
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
}

func NewBotter(s *state.State, ctx *context.Context) Botter {
	return Botter{
		BState:    s,
		Ctx:       *ctx,
		ComChan:   make(chan ComData),
		AllowList: []string{},
		awaitSong: true,
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
					println("reqnew")
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
					println("queue size too small..")
					continue
				}
				log.Println("request to play song")
				b.NewPlayData(b.Queue[b.QueueIndex])
				b.PlayData.Start()
				go b.PlayData.SendSongInfo()
			}
		}
	}
}
func (b *Botter) ClearQueue() {
	b.Queue = []QueueEntry{}
	b.QueueIndex = 0
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
