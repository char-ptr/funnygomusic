package commands

import (
	"fmt"
	"funnygomusic/bot"
	"funnygomusic/databaser"
	"github.com/diamondburned/arikawa/v3/gateway"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

func PlayCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	pathTo := ""
	query := strings.Join(args, " ")
	filepath.WalkDir("Y:/data/music", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if strings.Contains(strings.ToLower(info.Name()), strings.ToLower(query)) {
			pathTo = path
			return filepath.SkipAll
		}
		return nil
	})

	if pathTo == "" {
		_, err := b.BState.SendMessage(c.ChannelID, "unable to find, srry")
		if err != nil {
			log.Println(err)
			return
		}
		return

	}
	if !b.V.Open() {
		bot.JoinUsersVc(b, c.GuildID, c.Author.ID)
	}
	quelen := len(b.Queue.GetEntries())
	entry := databaser.NewIndexEntryFromPathDnc(pathTo)
	go b.Db.Create(&entry)
	b.Queue.Notify <- bot.NewPlaylistMessage(bot.PlaylistAdd).SetEntry(&entry)
	_, err := b.BState.SendMessage(c.ChannelID, fmt.Sprintf("added `%s` at index %d", entry.Title, quelen))
	if err != nil {
		log.Println(err)
		return
	}

}
