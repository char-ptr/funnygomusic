package commands

import (
	"fmt"
	"funnygomusic/bot"
	"funnygomusic/databaser"
	"funnygomusic/utils"
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
	if b.VoiceSes == nil {
		utils.JoinUsersVc(b, c)
	}
	quelen := len(b.Queue)
	entry := databaser.NewIndexEntryFromPathDnc(pathTo)
	b.Queue = append(b.Queue, entry)
	b.ComChan <- bot.NewItem
	_, err := b.BState.SendMessage(c.ChannelID, fmt.Sprintf("added `%s` at index %d", entry.Title, quelen))
	if err != nil {
		log.Println(err)
		return
	}

}
