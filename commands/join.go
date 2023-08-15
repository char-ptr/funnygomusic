package commands

import (
	"funnygomusic/bot"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func JoinCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	if b.VoiceSes.Open() {
		return
	}
	b.VoiceSes.JoinUsersVc(b, c.GuildID, c.Author.ID)
}
