package players

import (
	"context"
	"fmt"
	"funnygomusic/bot"
	"github.com/pkg/errors"
	"io"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type Youtube struct {
	url      string
	ctx      context.Context
	proc1    *exec.Cmd
	proc2    *exec.Cmd
	state    bot.PlayingState
	position int
	started  time.Time
	logger   *slog.Logger
	cmdout   io.Reader
}

func (yt *Youtube) Play(ctx context.Context) (io.Reader, error) {
	if yt.logger == nil {
		attr := slog.Attr{
			Key:   "scope",
			Value: slog.StringValue("Players/Youtube"),
		}
		yt.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}).WithAttrs([]slog.Attr{attr}))

	}
	cmd := exec.CommandContext(ctx, "yt-dlp",
		yt.url,
		"--downloader", "ffmpeg",
		"--force-keyframes-at-cuts",
		"--download-sections", fmt.Sprintf("*%d-", int(yt.position/1000)),
		"--postprocessor-args", "ffmpeg:\"-threads 1\"",
		"-x", "--audio-format", "opus",
		"--quiet", "-o", "-",
	)
	log.Println(cmd.String())
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stdout pipe")
	}

	yt.proc1 = cmd
	yt.ctx = ctx
	yt.state = bot.PSPlaying
	yt.started = time.Now()
	if yt.cmdout == nil {
		err = yt.proc1.Start()

	}

	cmd2 := exec.CommandContext(ctx, "ffmpeg", "-hide_banner", "-loglevel", "error",
		"-i", "pipe:",
		"-threads", "1",
		"-c:a", "libopus",
		"-b:a", "96k",
		"-frame_duration", strconv.Itoa(bot.FrameDuration),
		"-vbr", "off",
		"-f", "opus",
		"-",
	)
	cmd2.Stderr = os.Stderr

	cmd2.Stdin = stdout
	if yt.cmdout != nil {
		cmd2.Stdin = yt.cmdout
	}
	stdout, err = cmd2.StdoutPipe()
	yt.proc2 = cmd2
	cmd2.Start()

	if err != nil {
		yt.logger.Log(ctx, slog.LevelError, "Failed to start ffmpeg", "error", err)
		return nil, errors.Wrap(err, "failed to start ffmpeg")
	}
	return stdout, nil
}
func (yt *Youtube) Pause() error {
	if yt.state != bot.PSPlaying {
		yt.logger.Log(yt.ctx, slog.LevelError, "Cannot pause when not playing")
		return errors.New("cannot pause when not playing")
	}
	yt.position = int(time.Since(yt.started).Milliseconds())
	yt.state = bot.PSPaused
	err := yt.proc1.Cancel()
	err = yt.proc2.Cancel()
	if err != nil {
		yt.logger.Log(yt.ctx, slog.LevelError, "Failed to cancel ffmpeg", "error", err)
		return errors.Wrap(err, "failed to cancel ffmpeg")
	}
	return nil
}
func (yt *Youtube) Resume() (io.Reader, error) {
	return yt.Play(yt.ctx)
}
func (yt *Youtube) Seek(ms int) (io.Reader, error) {
	yt.position = ms
	if yt.state == bot.PSPlaying {
		err := yt.proc1.Cancel()
		err = yt.proc2.Cancel()
		if err != nil {
			yt.logger.Log(yt.ctx, slog.LevelError, "Failed to cancel ffmpeg", "error", err)
			return nil, errors.Wrap(err, "failed to cancel ffmpeg")
		}
	}
	return yt.Play(yt.ctx)
}
func (yt *Youtube) Stop() error {
	yt.state = bot.PSComplete
	err := yt.proc1.Cancel()
	err = yt.proc2.Cancel()
	if err != nil {
		yt.logger.Log(yt.ctx, slog.LevelError, "Failed to cancel ffmpeg", "error", err)
		return errors.Wrap(err, "failed to cancel ffmpeg")
	}
	return nil
}
func (yt *Youtube) Position() int {
	addOn := 0
	if yt.state == bot.PSPlaying {
		addOn = int(time.Since(yt.started).Milliseconds())
	}
	return yt.position + addOn
}
func (yt *Youtube) PositionDuration() time.Duration {
	return time.Duration(yt.Position()) * time.Millisecond
}
func (yt *Youtube) State() bot.PlayingState {
	return yt.state
}
func (yt *Youtube) WaitTillEnd() error {
	if yt.proc1 == nil {
		yt.logger.Log(yt.ctx, slog.LevelError, "No process to wait for")
		return errors.New("no process to wait for")
	}
	err := yt.proc1.Wait()
	err = yt.proc2.Wait()
	if err != nil {
		yt.logger.Log(yt.ctx, slog.LevelError, "Failed to wait for ffmpeg", "error", err)
		return err
	}
	yt.state = bot.PSComplete

	return nil
}

func NewYoutubePlayer(url string) *Youtube {
	return &Youtube{
		url: url,
	}
}
