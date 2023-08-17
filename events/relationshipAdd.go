package events

import (
	"funnygomusic/bot"
	"github.com/diamondburned/arikawa/v3/gateway"
	"log"
)

func OnRelationshipAdd(c *gateway.RelationshipAddEvent, b *bot.Botter) {
	log.Println("relationship add event", c.Type)
	//if c.Type == discord.IncomingFriendRequest {
	//	err := b.Gateway().Send(b.Ctx, &gateway.RelationshipAddEvent{
	//		Relationship: discord.Relationship{
	//			Type:   discord.SentFriendRequest,
	//			User:   c.User,
	//			UserID: c.UserID,
	//		},
	//	})
	//	if err != nil {
	//		log.Println("unable to accept friend request: ", err)
	//		return
	//	}
	//	_, err = b.SendMessage(discord.ChannelID(c.UserID), "hi :)")
	//	if err != nil {
	//		log.Println("unable to send message: ", err)
	//		return
	//	}
	//}
}
