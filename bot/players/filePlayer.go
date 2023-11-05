package players

import (
	"context"
	"fmt"
	"funnygomusic/bot"
	"github.com/pkg/errors"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"io"
	"log"
	"log/slog"
	"os"
	"time"
)

type StopNotification = int

const (
	SuccessProc StopNotification = iota
	FailedProc
)

type File struct {
	f        string
	ctx      context.Context
	state    bot.PlayingState
	ff       *ffmpeg_go.Stream
	position int
	started  time.Time
	logger   *slog.Logger
	pwer     *io.PipeWriter
	waiter   chan StopNotification
}

func (f *File) Play(ctx context.Context) (io.Reader, error) {
	if f.logger == nil {
		attr := slog.Attr{
			Key:   "scope",
			Value: slog.StringValue("Players/File"),
		}
		f.logger = slog.New(slog.NewTextHandler(os.Stdout, nil).WithAttrs([]slog.Attr{attr}))

	}
	if _, err := os.Stat(f.f); os.IsNotExist(err) {
		slog.Log(ctx, slog.LevelError, "File does not exist", "error", err)
		return nil, errors.Wrap(err, "file does not exist")
	}
	//cmd := exec.CommandContext(ctx, "ffmpeg", "-hide_banner", "-loglevel", "error",
	//	"-i", f.f,
	//	"-threads", "1",
	//	"-ss", fmt.Sprintf("%dms", f.position),
	//	"-c:a", "libopus",
	//	"-b:a", "96k",
	//	"-frame_duration", strconv.Itoa(bot.FrameDuration),
	//	"-vbr", "off",
	//	"-f", "opus",
	//	"-",
	//)
	//cmd.Stderr = os.Stderr
	//stdout, err := cmd.StdoutPipe()
	//if err != nil {
	//	return nil, errors.Wrap(err, "failed to get stdout pipe")
	//}
	//f.proc = cmd
	//f.ctx = ctx
	f.state = bot.PSPlaying
	f.started = time.Now()
	//err = f.proc.Start()
	//if err != nil {
	//	f.logger.Log(ctx, slog.LevelError, "Failed to start ffmpeg", "error", err)
	//	return nil, errors.Wrap(err, "failed to start ffmpeg")
	//}
	//return stdout, nil
	pr, pw := io.Pipe()
	f.pwer = pw
	ffip := ffmpeg_go.Input(f.f).Output("pipe:", ffmpeg_go.KwArgs{
		"format":         "opus",
		"vbr":            "off",
		"frame_duration": bot.FrameDuration,
		"ss":             fmt.Sprintf("%dms", f.position),
		"c:a":            "libopus",
		"b:a":            "96k",
		//"filter:a":       "volume=30",
	}).WithOutput(pw)
	go func() {
		err := ffip.Run()
		log.Println("ffmpeg exited", err)
		if err != nil {
			f.logger.Log(ctx, slog.LevelError, "Failed to start ffmpeg", "error", err)
			pw.CloseWithError(err)
			f.waiter <- FailedProc
			return
		}
		f.state = bot.PSComplete
		pw.Close()
		f.waiter <- SuccessProc

	}()

	return pr, nil
}
func (f *File) Pause() error {
	if f.state != bot.PSPlaying {
		f.logger.Log(f.ctx, slog.LevelError, "Cannot pause when not playing")
		return errors.New("cannot pause when not playing")
	}
	f.position = int(time.Since(f.started).Milliseconds())
	f.state = bot.PSPaused
	f.pwer.Close()
	//f.ff.Context.Done()
	return nil
}
func (f *File) Resume() (io.Reader, error) {
	return f.Play(f.ctx)
}
func (f *File) Seek(ms int) (io.Reader, error) {
	f.position = ms
	if f.state == bot.PSPlaying {
		f.state = bot.PSRestart
		f.pwer.Close()
		//f.ff.Context.Done()
	}
	return f.Play(f.ctx)
}
func (f *File) Stop() error {
	f.state = bot.PSComplete
	f.pwer.Close()
	//f.ff.Context.Done()
	return nil
}
func (f *File) Position() int {
	addOn := 0
	if f.state == bot.PSPlaying {
		addOn = int(time.Since(f.started).Milliseconds())
	}
	return f.position + addOn
}
func (f *File) PositionDuration() time.Duration {
	return time.Duration(f.Position()) * time.Millisecond
}
func (f *File) State() bot.PlayingState {
	return f.state
}
func (f *File) WaitTillEnd() error {
	log.Println("waiting for ffmpeg to exit")
	typer := <-f.waiter
	log.Println("ffmpeg exited2", typer)
	if typer == FailedProc {
		return errors.New("failed to start ffmpeg")
	}
	return nil
}

func NewFilePlayer(f string) *File {
	return &File{
		f:      f,
		waiter: make(chan StopNotification),
	}
}
