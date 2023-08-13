package events

import (
	"funnygomusic/bot"
	"log"

	"github.com/diamondburned/arikawa/v3/gateway"
)

func OnReady(c *gateway.ReadyEvent, b *bot.Botter) {
	b.MyId = c.User.ID
	b.MyUsername = c.User.Username
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
