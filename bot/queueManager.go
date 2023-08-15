package bot

import (
	"context"
	"fmt"
	"log"
	"time"
)

type PlaylistCmd int

const (
	PlaylistAdd PlaylistCmd = iota
	PlaylistRemove
	PlaylistClear
	PlaylistRepeat
	PlaylistShuffle
	CurrentStop
	CurrentSkip
	CurrentPause
	CurrentResume
	CurrentRestart
	CurrentSeek
	CurrentStart
	SongEnded
	SongPipeEnd
	SongProcEnd
	Jump
)

type QueueMessage struct {
	cmd      PlaylistCmd
	index    int
	seek     int
	newEntry *QueueEntry
}

func (p *QueueMessage) SetIndex(idx int) *QueueMessage {
	p.index = idx
	return p
}
func (p *QueueMessage) SetSeek(seek int) *QueueMessage {
	p.seek = seek
	return p
}
func (p *QueueMessage) SetEntry(entry *QueueEntry) *QueueMessage {
	p.newEntry = entry
	return p
}
func NewPlaylistMessage(msg PlaylistCmd) *QueueMessage {
	return &QueueMessage{cmd: msg}
}

type QueueManager struct {
	b           *Botter
	playlist    []QueueEntry
	index       int
	playingData *PlayingData
	Notify      chan *QueueMessage
}

func NewQueueManager(b *Botter) *QueueManager {

	return &QueueManager{
		b:        b,
		playlist: []QueueEntry{},
		index:    0,
		Notify:   make(chan *QueueMessage, 10),
	}

}

func (p *QueueManager) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-p.Notify:
			println("got message:", msg.cmd)
			switch msg.cmd {
			case PlaylistAdd:
				{
					p.playlist = append(p.playlist, *msg.newEntry)
					p.Notify <- NewPlaylistMessage(CurrentStart)

				}
			case PlaylistRemove:
				{
					p.playlist = append(p.playlist[:msg.index], p.playlist[msg.index+1:]...)
				}
			case PlaylistClear:
				{
					p.index = 0
					p.playlist = []QueueEntry{}
				}
			case CurrentStop:
				{
					p.index = 0
					p.playlist = nil
					if p.playingData != nil {
						p.playingData.Stop()
						p.playingData = nil
					}
				}
			case CurrentSkip:
				{
					if p.playingData != nil {
						p.playingData.Stop()
					}
				}
			case CurrentPause:
				{
					if p.GetPlayingState() == PSPlaying {
						p.playingData.Pause()
					}
				}
			case CurrentResume:
				{
					if p.GetPlayingState() == PSPaused {
						p.playingData.Resume()
					}
				}
			case CurrentRestart:
				{
					if p.playingData != nil {
						p.playingData.Restart()
					}
				}
			case CurrentStart:
				{
					if p.index >= len(p.playlist) {
						p.playingData = nil
						log.Println("too long.")
						p.b.SendMessage(p.b.SubChan, "Queue has ended")

						continue
					}
					if p.playingData == nil || (p.playingData != nil && p.playingData.state == PSComplete) {
						log.Println("request to play song")
						curr := p.GetCurrentSong()
						p.NewPlayData(curr)
						p.playingData.Start()
						p.b.SendMessage(p.b.SubChan, fmt.Sprintf("Now Playing:\nName: `%s`\nArtist: `%s`\nAlbum: `%s`\nFor `%s`", curr.Title, curr.Artist, curr.Album, curr.DurationStr()))
						//go p.playingData.SendSongInfo()
					}

				}
			case CurrentSeek:
				{
					if p.GetPlayingState() == PSPlaying {
						p.playingData.Seek(uint64(msg.seek))
					}
				}
			case SongEnded:
				{
					if len(p.playlist) == 0 {
						continue
					}
					p.index++
					p.playingData.state = PSComplete
					p.Notify <- NewPlaylistMessage(CurrentStart)
				}
			case SongPipeEnd:
				{
					if p.GetPlayingState() == PSComplete {
						p.Notify <- NewPlaylistMessage(SongEnded)
					} else if p.GetPlayingState() == PSRestart {
						p.playingData.Start()
					}
				}
			case SongProcEnd:
				{
					if p.playingData == nil {
						continue
					}
					if p.GetPlayingState() == PSPlaying {
						p.playingData.state = PSComplete
					}
				}
			case Jump:
				{
					if p.GetPlayingState() == PSNotPlaying {
						p.index = msg.index
						p.Notify <- NewPlaylistMessage(CurrentStart)
					} else {
						p.index = msg.index - 1
						p.playingData.Stop()
					}

				}
			}
		}
	}
}
func (p *QueueManager) GetCurrentSong() *QueueEntry {
	if p.index >= len(p.playlist) {
		return nil
	}
	return &p.playlist[p.index]
}
func (p *QueueManager) GetCurrentSongTime() time.Duration {
	if p.playingData == nil {
		return time.Duration(0)
	}
	return p.playingData.GetPlayingTime()
}
func (p *QueueManager) GetEntries() []QueueEntry {
	return p.playlist
}
func (p *QueueManager) GetIndex() int {
	return p.index
}

func (p *QueueManager) NewPlayData(entry *QueueEntry) {
	pd := NewPlayingData(p.b, entry)
	p.playingData = pd
}
func (p *QueueManager) GetPlayingState() PlayingState {
	if p.playingData == nil {
		return PSNotPlaying
	}
	return p.playingData.state
}
