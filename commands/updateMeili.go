package commands

import (
	"funnygomusic/bot"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func UpdateMeiliCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	b.MeiliUpdate()
	b.SendMessage(c.ChannelID, "ok")
}
