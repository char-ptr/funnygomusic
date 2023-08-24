package commands

import (
	"fmt"
	"funnygomusic/bot"
	"funnygomusic/bot/entries"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/meilisearch/meilisearch-go"
	"log"
	"strings"
)

func PlayIndexedCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	query := strings.Join(args, " ")
	//var songData []databaser.IndexedSong
	rmsg, err := b.Meili.Index("songs").Search(query, &meilisearch.SearchRequest{Limit: 1, Query: query})
	if err != nil || rmsg.Hits == nil || len(rmsg.Hits) == 0 {
		log.Println(err)
		b.SendMessage(c.ChannelID, "unable to find, srry")
		return
	}
	if !b.VoiceSes.Open() {
		b.VoiceSes.JoinUsersVc(b, c.GuildID, c.Author.ID)
	}
	var hit = rmsg.Hits[0].(map[string]interface{})

	fullRes := entries.Indexed{
		F:        hit["File"].(string),
		Title:    hit["Title"].(string),
		Artist:   hit["Artist"].(string),
		Album:    hit["Album"].(string),
		Duration: hit["Length"].(float64),
	}
	quelen := len(b.Queue.GetEntries())
	b.Queue.Notify <- bot.NewPlaylistMessage(bot.PlaylistAdd).SetEntry(&fullRes)
	b.SendMessage(c.ChannelID, fmt.Sprintf("added `%s (%s)` at index %d", fullRes.GetTitle(), fullRes.GetArtist(), quelen))
}
