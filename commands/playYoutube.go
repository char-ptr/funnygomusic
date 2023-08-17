package commands

import (
	"fmt"
	"funnygomusic/bot"
	"funnygomusic/bot/entries"
	"github.com/diamondburned/arikawa/v3/gateway"
	"strings"
)

func PlayYoutubeCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	query := strings.Join(args, " ")
	if !b.VoiceSes.Open() {
		b.VoiceSes.JoinUsersVc(b, c.GuildID, c.Author.ID)
	}
	quelen := len(b.Queue.GetEntries())
	yter, err := entries.NewUrl(query)
	if err != nil {
		b.SendMessage(c.ChannelID, "unable to find, srry")
		return
	}
	for _, vids := range yter {
		hold := vids
		b.Queue.Notify <- bot.NewPlaylistMessage(bot.PlaylistAdd).SetEntry(&hold)
	}
	if len(yter) == 1 {
		b.SendMessage(c.ChannelID, fmt.Sprintf("added `%s (%s)` at index %d", yter[0].GetTitle(), yter[0].GetArtist(), quelen))
		return
	}
	b.SendMessage(c.ChannelID, fmt.Sprintf("added `%d` at index %d", len(yter), quelen))
}
