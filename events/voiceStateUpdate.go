package events

import (
	"funnygomusic/bot"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func OnVoiceStateUpdate(c *gateway.VoiceStateUpdateEvent, b *bot.Botter) {
	if u := b.V.GetUser(c.UserID); u != nil {
		if !c.ChannelID.IsValid() {
			// user left
			b.V.DeleteUser(c.UserID)
		} else {
			vsu := bot.NewVoiceSessionUser(&c.VoiceState, b)
			b.V.UpdateUser(c.UserID, &vsu)
		}
	} else {
		vsu := bot.NewVoiceSessionUser(&c.VoiceState, b)
		b.V.AddUser(&vsu)
	}
}
