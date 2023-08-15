package commands

import (
	"funnygomusic/bot"
	"funnygomusic/databaser"
	"github.com/diamondburned/arikawa/v3/gateway"
	"strconv"
)

func AllowCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	for _, mem := range args {
		if memid, err := strconv.Atoi(mem); err == nil {
			b.AllowList = append(b.AllowList, mem)
			go b.Db.Create(&databaser.AllowedUser{ID: uint64(memid)})
		}
	}
}
