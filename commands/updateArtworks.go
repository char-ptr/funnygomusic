package commands

import (
	"funnygomusic/bot"
	"funnygomusic/databaser"
	"github.com/diamondburned/arikawa/v3/gateway"
	"os"
)

func UpdateArtworksCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	botOwner, exist := os.LookupEnv("BOT_OWNER")
	if !exist {
		return
	}
	if c.Author.ID.String() != botOwner {
		return
	}
	databaser.UpdateIndexedArtworks(b.Db)
	b.SendMessage(c.ChannelID, "ok")
}
