package commands

import (
	"funnygomusic/bot"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func PauseCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	if b.PlayData.Playing {
		b.PlayData.Pause()
	}
	b.PlayData.Paused = true
}
