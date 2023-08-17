package events

import (
	"funnygomusic/bot"
	"log"

	"github.com/diamondburned/arikawa/v3/gateway"
)

func OnReady(c *gateway.ReadyEvent, b *bot.Botter) {
	b.MyId = c.User.ID
	b.MyUsername = c.User.Username
	log.Println("connected as: ", c.User.DisplayOrUsername())
}
