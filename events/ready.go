package events

import (
	"funnygomusic/bot"
	"log"

	"github.com/diamondburned/arikawa/v3/gateway"
)

func onReady(c *gateway.ReadyEvent, le_bot *bot.Botter) {
	le_bot.MyId = c.User.ID
	// le_bot.BState.Gateway().Send(le_bot.Ctx, &gateway.UpdatePresenceCommand{
	// 	Status: discord.IdleStatus,
	// 	Activities: []discord.Activity{
	// 		{
	// 			Type:  discord.CustomActivity,
	// 			Emoji: &discord.Emoji{ID: discord.NullEmojiID, Name: "🤤"},
	// 			State: "Kafka...",
	// 		},
	// 	},
	// })
	log.Println("connected as: ", c.User.DisplayOrUsername())
}

func OnReady(le_bot *bot.Botter) func(c *gateway.ReadyEvent) {
	return func(c *gateway.ReadyEvent) {
		onReady(c, le_bot)
	}
}
