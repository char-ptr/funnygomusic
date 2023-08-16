package commands

import (
	"fmt"
	"funnygomusic/bot"
	"funnygomusic/bot/entries"
	"funnygomusic/databaser"
	"github.com/diamondburned/arikawa/v3/gateway"
	"strings"
)

func PlayIndexedCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	query := strings.Join(args, " ")
	res := databaser.TryFindSong(query, b.Db)
	if res == nil {
		b.SendMessage(c.ChannelID, "unable to find, srry")
		return
	}
	if !b.VoiceSes.Open() {
		b.VoiceSes.JoinUsersVc(b, c.GuildID, c.Author.ID)
	}
	fullRes := entries.NewIndexedFromDb(res.ID.String(), b.Db)

	quelen := len(b.Queue.GetEntries())
	b.Queue.Notify <- bot.NewPlaylistMessage(bot.PlaylistAdd).SetEntry(&fullRes)
	b.SendMessage(c.ChannelID, fmt.Sprintf("added `%s (%s)` at index %d", fullRes.GetTitle(), fullRes.GetArtist(), quelen))
}
