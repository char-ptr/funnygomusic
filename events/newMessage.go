package events

import (
	"funnygomusic/bot"
	"funnygomusic/commands"
	"strings"

	"github.com/diamondburned/arikawa/v3/gateway"
	"golang.org/x/exp/slices"
)

var (
	Commands = map[string]commands.Command{
		"join":      commands.JoinCommand,
		"play":      commands.PlayIndexedCommand,
		"p":         commands.PlayIndexedCommand,
		"playf":     commands.PlayFileCommand,
		"pf":        commands.PlayFileCommand,
		"py":        commands.PlayYoutubeCommand,
		"pause":     commands.PauseCommand,
		"resume":    commands.ResumeCommand,
		"leave":     commands.LeaveCommand,
		"fuckoff":   commands.LeaveCommand,
		"skip":      commands.SkipCommand,
		"queue":     commands.QueueCommand,
		"q":         commands.QueueCommand,
		"seek":      commands.SeekCommand,
		"jump":      commands.JumpCommand,
		"song-info": commands.SongInfoCommand,
		"allow":     commands.AllowCommand,
		"restart":   commands.RestartCommand,
		"testvce":   commands.TestVceCommand,
		"idxdir":    commands.IndexDirCommand,
	}
)

func OnMessage(c *gateway.MessageCreateEvent, b *bot.Botter) {
	if !strings.HasPrefix(c.Message.Content, "`") || !(slices.Contains(b.AllowList, c.Author.ID.String())) {
		return
	}
	commandArgs := strings.Split(c.Message.Content, " ")
	command := commandArgs[0][1:]
	commandArgs = commandArgs[1:]
	b.SubChan = c.ChannelID
	for k, v := range Commands {
		if k == command {
			v(c, b, commandArgs)
			return
		}
	}
}
