package bot

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/dhowden/tag"
	"github.com/diamondburned/arikawa/v3/voice/voicegateway"
	"github.com/diamondburned/oggreader"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

const (
	FrameDuration = 60 // ms
	TimeIncrement = 2880
)

type PlayingData struct {
	bot           *Botter
	cmd           *exec.Cmd
	Playing       bool
	Paused        bool
	Complete      bool
	CustomFilters string
	Bitrate       int
	time_played   uint64
	time_started  time.Time
	file          string
	outstm        io.ReadCloser
}

func NewPlayingData(bot *Botter, file string) *PlayingData {
	pd := &PlayingData{
		bot:          bot,
		cmd:          nil,
		file:         file,
		Bitrate:      94,
		time_played:  0,
		time_started: time.Now(),
	}
	return pd
}
func (p *PlayingData) Stop() {

	p.outstm.Close()
	p.cmd.Cancel()
	p.Paused = false
	p.Playing = false
	p.Complete = true
	p.time_played = 0
	p.time_started = time.Now()
	p.bot.VoiceSes.Speaking(p.bot.Ctx, voicegateway.NotSpeaking)
}
func (p *PlayingData) Start() error {
	p.RestoreCmd()
	err := p.cmd.Start()
	if err != nil {
		log.Println(errors.Wrap(err, "failed to start ffmpeg"))
	}

	p.bot.VoiceSes.Speaking(p.bot.Ctx, voicegateway.Microphone)

	p.Playing = true
	p.Paused = false
	go p.PipeIntoStream()
	go p.playLoop()
	return nil
}
func (p *PlayingData) playLoop() {
	p.cmd.Wait()
	if !p.Paused {
		p.Complete = true
	}
}
func (p *PlayingData) PipeIntoStream() error {
	if err := oggreader.DecodeBuffered(p.bot.VoiceSes, p.outstm); err != nil {
		log.Println(errors.Wrap(err, "failed to decode ogg"))
	}
	println("finished piping")
	if p.Complete {
		println("finished playing")
		p.Stop()
		p.bot.ComChan <- SongEnded
	} else {
		p.Start()
	}
	return nil
}

func (p *PlayingData) Restart() {
	// set p.time_played to distance from now and p.time_started

	p.time_played += uint64(time.Since(p.time_started).Milliseconds())
	p.time_started = time.Now()
	p.Paused = true
	p.cmd.Cancel()

}
func (p *PlayingData) Pause() {
	p.time_played += uint64(time.Since(p.time_started).Milliseconds())
	p.bot.VoiceSes.Speaking(p.bot.Ctx, voicegateway.NotSpeaking)
	defer p.outstm.Close()
	p.cmd.Cancel()
	p.time_started = time.Now()
	p.Playing = false
	p.Paused = true

}

func (p *PlayingData) Resume() {

	p.Playing = true
	p.Paused = false
	p.Start()
	p.time_started = time.Now()

}
func (p *PlayingData) Seek(timer uint64) {

	p.time_played = timer
	p.time_started = time.Now()
	p.Restart()
}
func (p *PlayingData) RestoreCmd() error {
	temp_cmd := exec.CommandContext(p.bot.Ctx,
		"ffmpeg", "-hide_banner", "-loglevel", "error",
		// Streaming is slow, so a single thread is all we need.
		"-threads", "1",
		"-ss", fmt.Sprintf("%dms", p.time_played),
		// Input file.
		"-i", p.file,
		// Output format; leave as "libopus".
		"-c:a", "libopus",
		// Bitrate in kilobits. This doesn't matter, but I recommend 96k as the
		// sweet spot.
		"-b:a", fmt.Sprintf("%dk", p.Bitrate),

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
	if p.CustomFilters != "" {
		log.Println("adding custom filters:" + p.CustomFilters)
		temp_cmd.Args = slices.Insert(temp_cmd.Args, 10, "-af", p.CustomFilters)
	}
	temp_cmd.Stderr = os.Stderr
	stdout, err := temp_cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "failed to get stdout pipe")
	}
	p.outstm = stdout
	p.cmd = temp_cmd
	return nil
}
func (p *PlayingData) SendSongInfo() {
	thef, _ := os.Open(p.bot.CurrentPlayingSong().Path)
	defer thef.Close()
	tag, err := tag.ReadFrom(thef)
	if err != nil {
		p.bot.BState.SendMessage(p.bot.SubChan, "error opening file")
		return
	}
	msg_cnt := fmt.Sprintf("Currently playing:\nName: `%s`\nArtist: `%s`\nAlbum: `%s`", tag.Title(), tag.Artist(), tag.Album())

	p.bot.BState.SendMessage(p.bot.SubChan, msg_cnt)
}
