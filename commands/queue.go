package commands

import (
	"fmt"
	"funnygomusic/bot"
	"funnygomusic/utils"
	"github.com/diamondburned/arikawa/v3/gateway"
	"log"
	"strconv"
	"strings"
)

type Modes int

const (
	List Modes = iota
	Remove
	Clear
)

func ClearQueue(b *bot.Botter) {
	b.Queue.Notify <- bot.NewPlaylistMessage(bot.PlaylistClear)
	//b.BState.SendMessage(b.Ctx, "cleared queue")
}

func RemoveEntry(c *gateway.MessageCreateEvent, b *bot.Botter, idx int) {
	tmpQueue := b.Queue.GetEntries()
	if idx > len(tmpQueue) {
		b.BState.SendMessage(c.ChannelID, "out of bounds")
		return
	}
	whatsThere := tmpQueue[idx]
	b.Queue.Notify <- bot.NewPlaylistMessage(bot.PlaylistRemove).SetIndex(idx)
	b.BState.SendMessage(c.ChannelID, fmt.Sprintf("removed `%s`", whatsThere.Title))
}

func ListQueue(c *gateway.MessageCreateEvent, b *bot.Botter) {
	tmpQueue := b.Queue.GetEntries()
	if len(tmpQueue) == 0 {
		b.BState.SendMessage(c.ChannelID, "nothing in queue")
		return
	}
	msgCnt := "queue:```"
	for k, v := range tmpQueue {
		msgCnt += fmt.Sprintf("%d. %s - %s\n", k, v.Artist, v.Title)
	}
	msgCnt += "```"
	_, err := b.BState.SendMessage(c.ChannelID, msgCnt)
	if err != nil {
		log.Println(err)
	}

}

func QueueCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	if len(args) == 0 {
		ListQueue(c, b)
		return
	}
	switch strings.ToLower(args[0]) {
	case "list":
		{
			ListQueue(c, b)
		}
	case "remove":
		{
			if len(args) < 2 {
				b.BState.SendMessage(c.ChannelID, "invalid arg length")
				return
			}
			idx, err := strconv.Atoi(utils.GetIndex(args, 1))
			if err != nil {
				b.BState.SendMessage(c.ChannelID, "unable to parse index as int")
				return
			}
			RemoveEntry(c, b, idx)
		}
	case "clear":
		{
			ClearQueue(b)
		}
	}
}
