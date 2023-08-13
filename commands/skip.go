package commands

import (
	"funnygomusic/bot"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func SkipCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	b.Queue.Notify <- bot.NewPlaylistMessage(bot.CurrentSkip)
}
