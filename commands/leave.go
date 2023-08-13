package commands

import (
	"funnygomusic/bot"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func LeaveCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	if b.PlayData != nil {
		b.PlayData.Stop()
		b.PlayData = nil
	}
	b.ClearQueue()
	b.VoiceSes.Leave(b.Ctx)
	b.VoiceSes = nil
}
