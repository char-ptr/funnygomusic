package events

import (
	"funnygomusic/bot"

	"github.com/diamondburned/arikawa/v3/gateway"
)

func OnVoiceStateUpdate(c *gateway.VoiceStateUpdateEvent, b *bot.Botter) {
	if !b.VoiceSes.Open() {
		return
	}
	if u := b.VoiceSes.GetUser(c.UserID); u != nil {
		if !c.ChannelID.IsValid() {
			// user left
			b.VoiceSes.DeleteUser(c.UserID)
		} else {
			vsu := bot.NewVoiceSessionUser(&c.VoiceState, b)
			if vsu.ID == b.MyId {
				if vsu.Muted || c.SelfDeaf {
					b.Queue.Notify <- bot.NewPlaylistMessage(bot.CurrentPause)
				} else {
					b.Queue.Notify <- bot.NewPlaylistMessage(bot.CurrentResume)
				}
			}
			b.VoiceSes.UpdateUser(c.UserID, &vsu)
		}
	} else {
		vsu := bot.NewVoiceSessionUser(&c.VoiceState, b)
		b.VoiceSes.AddUser(&vsu)
	}
}
