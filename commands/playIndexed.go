package commands

import (
	"fmt"
	"funnygomusic/bot"
	"funnygomusic/bot/entries"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/hako/durafmt"
	"github.com/meilisearch/meilisearch-go"
	"log"
	"net/url"
	"strings"
)

func PlayIndexedCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	query := strings.Join(args, " ")
	//var songData []databaser.IndexedSong
	quelen := len(b.Queue.GetEntries())
	queuedur := b.Queue.GetDuration()
	urlMaybem, ok := url.Parse(query)
	var entry []bot.QueueEntry
	if ok == nil && (urlMaybem.Scheme == "https" || urlMaybem.Scheme == "http") {
		entriesb, err := entries.NewUrl(query)
		if err != nil {
			b.SendMessage(c.ChannelID, "unable to find, srry2")
			return
		}
		for _, ent := range entriesb {
			holder := ent
			entry = append(entry, &holder)
		}

	} else {
		log.Println("unable to parse url")
		rmsg, err := b.Meili.Index("songs").Search(query, &meilisearch.SearchRequest{Limit: 1, Query: query})
		if err != nil || rmsg.Hits == nil || len(rmsg.Hits) == 0 {
			log.Println(err)
			b.SendMessage(c.ChannelID, "unable to find, srry1")
			return
		}
		var hit = rmsg.Hits[0].(map[string]interface{})
		entryu := &entries.Indexed{
			F:        hit["File"].(string),
			Title:    hit["Title"].(string),
			Artist:   hit["Artist"].(string),
			Album:    hit["Album"].(string),
			Duration: hit["Length"].(float64),
		}
		entry = append(entry, entryu)

	}
	log.Println("add entries = ", entry)
	if len(entry) == 0 {
		b.SendMessage(c.ChannelID, "unable to find, srry3")
		return
	}
	if !b.VoiceSes.Open() {
		b.VoiceSes.JoinUsersVc(b, c.GuildID, c.Author.ID)
	}
	for _, vids := range entry {
		hold := vids
		b.Queue.Notify <- bot.NewPlaylistMessage(bot.PlaylistAdd).SetEntry(hold)
	}
	b.SendMessage(c.ChannelID, fmt.Sprintf("added `%s (%s)` at index %d. will play in ~%s", entry[0].GetTitle(), entry[0].GetArtist(), quelen, durafmt.Parse(queuedur).LimitFirstN(2)))
}
