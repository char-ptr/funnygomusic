package commands

import (
	"funnygomusic/bot"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func JoinCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	if b.V.Open() {
		return
	}
	bot.JoinUsersVc(b, c.GuildID, c.Author.ID)
}
