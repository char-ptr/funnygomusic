package commands

import (
	"fmt"
	"funnygomusic/bot"
	"funnygomusic/utils"
	"github.com/diamondburned/arikawa/v3/gateway"
	"golang.org/x/exp/slices"
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
	if b.PlayData != nil {
		b.PlayData.Stop()
		b.PlayData = nil
	}
	b.ClearQueue()
}

func RemoveEntry(c *gateway.MessageCreateEvent, b *bot.Botter, idx int) {
	if idx > len(b.Queue) {
		b.BState.SendMessage(c.ChannelID, "out of bounds")
		return
	}
	whatsThere := b.Queue[idx]
	b.Queue = slices.Delete(b.Queue, idx, idx+1)
	b.BState.SendMessage(c.ChannelID, fmt.Sprintf("removed `%s`", whatsThere.Title))
}

func ListQueue(c *gateway.MessageCreateEvent, b *bot.Botter) {
	if len(b.Queue) == 0 {
		b.BState.SendMessage(c.ChannelID, "nothing in queue")
		return
	}
	msgCnt := "queue:```"
	for k, v := range b.Queue {
		msgCnt += fmt.Sprintf("%d. %s - %s\n", k, v.Artist, v.Title)
	}
	msgCnt += "```"
	b.BState.SendMessage(c.ChannelID, msgCnt)

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
