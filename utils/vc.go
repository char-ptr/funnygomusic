package utils

import (
	"funnygomusic/bot"
	"time"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/voice"
	"github.com/diamondburned/arikawa/v3/voice/udp"
	"github.com/pkg/errors"
)

func JoinUsersVc(Bot *bot.Botter, c *gateway.MessageCreateEvent) error {
	vs, err := voice.NewSession(Bot.BState)
	if err != nil {
		return errors.Wrap(err, "cannot make new voice session")
	}
	Bot.VoiceSes = vs
	vs.SetUDPDialer(udp.DialFuncWithFrequency(
		bot.FrameDuration*time.Millisecond, // correspond to -frame_duration
		bot.TimeIncrement,
	))
	uservs, err := Bot.BState.VoiceState(c.GuildID, c.Author.ID)
	if err != nil {
		return errors.Wrap(err, "cannot get voice state")
	}
	vs.JoinChannel(Bot.Ctx, uservs.ChannelID, false, false)
	return nil
}
