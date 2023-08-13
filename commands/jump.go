package commands

import (
	"funnygomusic/bot"
	"funnygomusic/utils"
	"github.com/diamondburned/arikawa/v3/gateway"
	"strconv"
)

func JumpCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	idx, err := strconv.Atoi(utils.GetIndex(args, 0))
	if err != nil {
		b.BState.SendMessage(c.ChannelID, "what")
		return
	}
	if idx > len(b.Queue.GetEntries()) {
		b.BState.SendMessage(c.ChannelID, "out of bounds")
		return
	}
	b.Queue.Notify <- bot.NewPlaylistMessage(bot.Jump).SetIndex(idx)

}
