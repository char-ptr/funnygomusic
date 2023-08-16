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
	yter := entries.NewUrl(query)
	b.Queue.Notify <- bot.NewPlaylistMessage(bot.PlaylistAdd).SetEntry(&yter)
	b.SendMessage(c.ChannelID, fmt.Sprintf("added `%s (%s)` at index %d", yter.GetTitle(), yter.GetArtist(), quelen))
}
