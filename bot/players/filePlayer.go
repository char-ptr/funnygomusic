package players

import (
	"context"
	"fmt"
	"funnygomusic/bot"
	"github.com/pkg/errors"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type StopNotification = int

const (
	SuccessProc = iota
	FailedProc
)

type File struct {
	f        string
	ctx      context.Context
	proc     *exec.Cmd
	state    bot.PlayingState
	position int
	started  time.Time
	logger   *slog.Logger
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
	cmd := exec.CommandContext(ctx, "ffmpeg", "-hide_banner", "-loglevel", "error",
		"-i", f.f,
		"-threads", "1",
		"-ss", fmt.Sprintf("%dms", f.position),
		"-c:a", "libopus",
		"-b:a", "96k",
		"-frame_duration", strconv.Itoa(bot.FrameDuration),
		"-vbr", "off",
		"-f", "opus",
		"-",
	)
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stdout pipe")
	}
	f.proc = cmd
	f.ctx = ctx
	f.state = bot.PSPlaying
	f.started = time.Now()
	err = f.proc.Start()
	if err != nil {
		f.logger.Log(ctx, slog.LevelError, "Failed to start ffmpeg", "error", err)
		return nil, errors.Wrap(err, "failed to start ffmpeg")
	}
	return stdout, nil
}
func (f *File) Pause() error {
	if f.state != bot.PSPlaying {
		f.logger.Log(f.ctx, slog.LevelError, "Cannot pause when not playing")
		return errors.New("cannot pause when not playing")
	}
	f.position = int(time.Since(f.started).Milliseconds())
	f.state = bot.PSPaused
	err := f.proc.Cancel()
	if err != nil {
		f.logger.Log(f.ctx, slog.LevelError, "Failed to cancel ffmpeg", "error", err)
		return errors.Wrap(err, "failed to cancel ffmpeg")
	}
	return nil
}
func (f *File) Resume() (io.Reader, error) {
	return f.Play(f.ctx)
}
func (f *File) Seek(ms int) (io.Reader, error) {
	f.position = ms
	if f.state == bot.PSPlaying {
		err := f.proc.Cancel()
		if err != nil {
			f.logger.Log(f.ctx, slog.LevelError, "Failed to cancel ffmpeg", "error", err)
			return nil, errors.Wrap(err, "failed to cancel ffmpeg")
		}
	}
	return f.Play(f.ctx)
}
func (f *File) Stop() error {
	f.state = bot.PSComplete
	err := f.proc.Cancel()
	if err != nil {
		f.logger.Log(f.ctx, slog.LevelError, "Failed to cancel ffmpeg", "error", err)
		return errors.Wrap(err, "failed to cancel ffmpeg")
	}
	return nil
}
func (f *File) Position() int {
	return f.position
}
func (f *File) PositionDuration() time.Duration {
	return time.Duration(f.position) * time.Millisecond
}
func (f *File) State() bot.PlayingState {
	return f.state
}
func (f *File) WaitTillEnd() error {
	if f.proc == nil {
		f.logger.Log(f.ctx, slog.LevelError, "No process to wait for")
		return errors.New("no process to wait for")
	}
	err := f.proc.Wait()
	if err != nil {
		f.logger.Log(f.ctx, slog.LevelError, "Failed to wait for ffmpeg", "error", err)
		return err
	}
	f.state = bot.PSComplete

	return nil
}

func NewFilePlayer(f string) *File {
	return &File{
		f: f,
	}
}
