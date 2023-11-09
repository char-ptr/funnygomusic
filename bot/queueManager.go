package bot

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/diamondburned/oggreader"
	"github.com/hako/durafmt"
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
	newEntry QueueEntry
	cmd      PlaylistCmd
	index    int
	seek     int
}

func (p *QueueMessage) SetIndex(idx int) *QueueMessage {
	p.index = idx
	return p
}

func (p *QueueMessage) SetSeek(seek int) *QueueMessage {
	p.seek = seek
	return p
}

func (p *QueueMessage) SetEntry(entry QueueEntry) *QueueMessage {
	p.newEntry = entry
	return p
}

func NewPlaylistMessage(msg PlaylistCmd) *QueueMessage {
	return &QueueMessage{cmd: msg}
}

type QueueManager struct {
	player   Player
	b        *Botter
	Notify   chan *QueueMessage
	logger   *slog.Logger
	playlist []QueueEntry
	index    int
}

func NewQueueManager(b *Botter) *QueueManager {
	attr := slog.String("scope", "Bot/QueueManager")
	txtHndlr := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}).WithAttrs([]slog.Attr{attr})
	logger := slog.New(txtHndlr)
	return &QueueManager{
		b:        b,
		playlist: []QueueEntry{},
		index:    0,
		Notify:   make(chan *QueueMessage, 10),
		logger:   logger,
	}
}

func (qm *QueueManager) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-qm.Notify:
			switch msg.cmd {
			case PlaylistAdd:
				{
					qm.logger.Debug("adding song to queue", "song", msg.newEntry.GetTitle())
					qm.playlist = append(qm.playlist, msg.newEntry)
					go func() {
						qm.Notify <- NewPlaylistMessage(CurrentStart)
					}()

				}
			case PlaylistRemove:
				{
					qm.playlist = append(qm.playlist[:msg.index], qm.playlist[msg.index+1:]...)
				}
			case PlaylistClear:
				{
					qm.index = 0
					qm.playlist = []QueueEntry{}
				}
			case CurrentStop:
				{
					qm.index = 0
					qm.playlist = nil
					if qm.player != nil {
						qm.player.Stop()
						qm.player = nil
					}
				}
			case CurrentSkip:
				{
					if qm.player != nil {
						qm.player.Stop()
					}
				}
			case CurrentPause:
				{
					if qm.GetPlayingState() == PSPlaying {
						qm.player.Pause()
						qm.b.VoiceSes.Speaking(false, ctx)
					}
				}
			case CurrentResume:
				{
					if qm.GetPlayingState() == PSPaused {
						r, _ := qm.player.Resume()
						qm.b.VoiceSes.Speaking(true, ctx)
						go qm.WriteData(r)
						go qm.AlertUponEnd()
					}
				}
			case CurrentStart:
				{
					if qm.index >= len(qm.playlist) {
						qm.player = nil
						qm.logger.Debug("queue has ended", "size", len(qm.playlist))
						qm.b.VoiceSes.Speaking(false, ctx)
						qm.b.SendMessage(qm.b.SubChan, "Queue has ended")

						continue
					}
					if qm.player == nil || qm.GetPlayingState() == PSComplete {
						curr := *qm.GetCurrentSong()
						qm.logger.Debug("request to play song", "index", qm.index, "song", curr)
						qm.player = curr.GetPlayer()
						qm.b.VoiceSes.Speaking(true, ctx)
						r, e := qm.player.Play(ctx)
						if e != nil {
							qm.logger.Debug("failed to play song, removing song", "error", e)
							go func() {
								qm.Notify <- NewPlaylistMessage(PlaylistRemove).SetIndex(qm.index)
								qm.Notify <- NewPlaylistMessage(CurrentStart)
							}()
							continue
						}
						go qm.WriteData(r)
						qm.b.SendMessage(qm.b.SubChan, fmt.Sprintf(">>> **Nоw Plаying**:** **%[5]s `%[1]s`\nArtіst:** **%[5]s `%[2]s`\n%[5]sAlbum:** **%[5]s `%[3]s`%[5]s\nFоr `%[4]s`",
							curr.GetTitle(), curr.GetArtist(), curr.GetAlbum(), durafmt.Parse(time.Duration(curr.GetDuration())*time.Millisecond).LimitFirstN(2), strings.Repeat("​", randInt(0, 5))))
						go qm.AlertUponEnd()
					}

				}
			case CurrentSeek:
				{
					if qm.GetPlayingState() == PSPlaying {
						r, _ := qm.player.Seek(msg.seek)
						go qm.WriteData(r)
						go qm.AlertUponEnd()
					}
				}
			case SongEnded:
				{
					if len(qm.playlist) == 0 {
						continue
					}
					qm.index++
					go func() {
						qm.Notify <- NewPlaylistMessage(CurrentStart)
					}()
				}
			case Jump:
				{
					if qm.GetPlayingState() == PSNotPlaying {
						qm.index = msg.index
						go func() {
							qm.Notify <- NewPlaylistMessage(CurrentStart)
						}()
					} else {
						qm.index = msg.index - 1
						qm.player.Stop()
					}
				}
			}
		}
	}
}

func (qm *QueueManager) AlertUponEnd() {
	err := qm.player.WaitTillEnd()
	if err != nil {
		return
	}
	if qm.GetPlayingState() == PSComplete {
		qm.logger.Debug("Song ended, complete")
	} else {
		qm.logger.Debug("Song ended, but not complete")
	}
}

func (qm *QueueManager) WriteData(reader io.Reader) {
	qm.logger.Debug("piping data -> ogg -> discord voice")
	if err := oggreader.DecodeBuffered(qm.b.VoiceSes.GetSession(), reader); err != nil {
		qm.logger.Log(qm.b.Ctx, slog.LevelError, "Failed to decode ogg", "error", err)
	}
	//_, err := io.Copy(qm.b.VoiceSes.GetSession(), reader)
	//if err != nil {
	//	qm.logger.Log(qm.b.Ctx, slog.LevelError, "Failed to copy ogg", "error", err)
	//	return
	//}
	if qm.GetPlayingState() == PSComplete {
		qm.Notify <- NewPlaylistMessage(SongEnded)
	} else {
		qm.logger.Debug("Song ended, but not complete(2)")
	}
}

func (qm *QueueManager) GetCurrentSong() *QueueEntry {
	if qm.index >= len(qm.playlist) {
		return nil
	}
	return &qm.playlist[qm.index]
}

func (qm *QueueManager) GetCurrentSongTime() time.Duration {
	if qm.player == nil {
		return time.Duration(0)
	}
	return qm.player.PositionDuration()
}

func (qm *QueueManager) GetEntries() []QueueEntry {
	return qm.playlist
}

func (qm *QueueManager) GetIndex() int {
	return qm.index
}

func (qm *QueueManager) GetDuration() (dur time.Duration) {
	for _, entry := range qm.playlist[qm.index:] {
		dur += time.Duration(entry.GetDuration()) * time.Millisecond
	}
	return
}

func (qm *QueueManager) GetPlayingState() PlayingState {
	if qm.player == nil {
		return PSNotPlaying
	}
	return qm.player.State()
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

//func (qm *QueueManager) SaveQueue(db *gorm.DB) {
//	var entries []string
//	for _, entry := range qm.playlist {
//		entries = append(entries, entry.GetID())
//	}
//	db.Model().Create(&entries)
//}
