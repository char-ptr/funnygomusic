package bot

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/diamondburned/oggreader"
	"github.com/pkg/errors"
)

type PlayingState int

const (
	PSPlaying PlayingState = iota
	PSPaused
	PSComplete
	PSRestart
	PSNotPlaying
)

const (
	FrameDuration = 60 // ms
	TimeIncrement = 2880
)

type PlayingData struct {
	bot       *Botter
	cmd       *exec.Cmd
	state     PlayingState
	duration  time.Duration
	started   time.Time
	entry     *QueueEntry
	outStream io.ReadCloser
}

func NewPlayingData(bot *Botter, entry *QueueEntry) *PlayingData {
	pd := &PlayingData{
		bot:      bot,
		cmd:      nil,
		entry:    entry,
		duration: time.Duration(0),
		started:  time.Now(),
	}
	return pd
}
func (p *PlayingData) Stop() {

	p.state = PSComplete
	p.cmd.Cancel()
	p.bot.V.Speaking(false)

}
func (p *PlayingData) Start() error {
	p.RestoreCmd()
	err := p.cmd.Start()
	if err != nil {
		log.Println(errors.Wrap(err, "failed to start ffmpeg"))
		return err
	}

	p.bot.V.Speaking(true)

	p.state = PSPlaying
	p.started = time.Now()
	go p.PipeIntoStream()
	go p.playLoop()
	return nil
}
func (p *PlayingData) playLoop() {
	p.cmd.Wait()
	p.bot.Queue.Notify <- NewPlaylistMessage(SongProcEnd)
}
func (p *PlayingData) PipeIntoStream() {
	if err := oggreader.DecodeBuffered(p.bot.V.GetSession(), p.outStream); err != nil {
		log.Println("ogg reader errored")
	}
	p.bot.Queue.Notify <- NewPlaylistMessage(SongPipeEnd)

}

func (p *PlayingData) Restart() {
	// set p.time_played to distance from now and p.time_started

	p.duration = p.GetPlayingTime()
	p.state = PSRestart
	p.cmd.Cancel()

}
func (p *PlayingData) Pause() {
	p.duration = p.GetPlayingTime()
	p.state = PSPaused
	p.cmd.Cancel()
	p.bot.V.Speaking(false)

}

func (p *PlayingData) Resume() {
	p.Start()
}
func (p *PlayingData) Seek(seconds uint64) {

	p.duration = time.Duration(seconds) * time.Millisecond
	p.state = PSRestart
	p.cmd.Cancel()
}
func (p *PlayingData) RestoreCmd() error {
	tempCmd := exec.CommandContext(p.bot.Ctx,
		"ffmpeg", "-hide_banner", "-loglevel", "error",
		// Streaming is slow, so a single thread is all we need.
		"-threads", "1",
		"-ss", fmt.Sprintf("%dms", p.duration.Milliseconds()),
		// Input file.
		"-i", p.entry.Path,
		// Output format; leave as "libopus".
		"-c:a", "libopus",
		// Bitrate in kilobits. This doesn't matter, but I recommend 96k as the
		// sweet spot.
		"-b:a", "96k",

		// Frame duration should be the same as what's given into
		// udp.DialFuncWithFrequency.
		"-frame_duration", strconv.Itoa(FrameDuration),
		// Disable variable bitrate to keep packet sizes consistent. This is
		// optional.
		"-vbr", "off",
		// Output format, which is opus, so we need to unwrap the opus file.
		"-f", "opus",
		"-",
	)

	tempCmd.Stderr = os.Stderr
	stdout, err := tempCmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "failed to get stdout pipe")
	}
	p.outStream = stdout
	p.cmd = tempCmd
	return nil
}
func (p *PlayingData) GetPlayingTime() time.Duration {
	sincet := time.Since(p.started)
	if p.state != PSPlaying {
		sincet = time.Duration(0)
	}
	return sincet + p.duration
}
