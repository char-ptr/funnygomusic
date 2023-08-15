package commands

import (
	"funnygomusic/bot"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func ResumeCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	if b.Queue.GetPlayingState() == bot.PSPaused {
		b.Queue.Notify <- bot.NewPlaylistMessage(bot.CurrentResume)
	}
}
