package commands

import (
	"context"
	"funnygomusic/bot"
	"funnygomusic/databaser"
	"github.com/diamondburned/arikawa/v3/gateway"
	"os"
	"strings"
)

func IndexDirCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	botOwner, exist := os.LookupEnv("BOT_OWNER")
	if !exist {
		return
	}
	if c.Author.ID.String() != botOwner {
		return
	}
	argsJoined := strings.Join(args, " ")
	newIndexer := databaser.NewIndexer(b.Db)
	ctx2 := context.TODO()
	newIndexer.IndexDirectory(argsJoined, ctx2)
}
