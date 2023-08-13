package commands

import (
	"funnygomusic/bot"
	"funnygomusic/utils"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func JoinCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	if b.VoiceSes != nil {
		return
	}
	utils.JoinUsersVc(b, c)
}
