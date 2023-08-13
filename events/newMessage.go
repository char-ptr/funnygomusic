package events

import (
	"fmt"
	"funnygomusic/bot"
	"funnygomusic/databaser"
	"funnygomusic/utils"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/gateway"
	"golang.org/x/exp/slices"
)

func OnMessage(c *gateway.MessageCreateEvent, b *bot.Botter) {
	if !strings.HasPrefix(c.Message.Content, "`") || !(slices.Contains(b.AllowList, c.Author.ID.String())) {
		return
	}
	command_args := strings.Split(c.Message.Content, " ")
	command := command_args[0][1:]
	command_args = command_args[1:]
	b.SubChan = c.ChannelID
	switch command {
	case "join":
		{
			utils.JoinUsersVc(b, c)
		}

	case "play":
		{
			// if b.PlayData != nil && b.PlayData.Playing {
			// 	log.Println("already playing")
			// 	return
			// }
			path_to := strings.Join(command_args, " ")
			if _, err := os.Stat(path_to); os.IsNotExist(err) {
				log.Println("file does not exist")
				return
			}
			if b.VoiceSes == nil {
				utils.JoinUsersVc(b, c)
			}
			quelen := len(b.Queue)
			b.Queue = append(b.Queue, databaser.NewIndexEntryFromPathDnc(path_to))
			b.ComChan <- bot.NewItem
			_, errm := b.BState.SendMessage(c.ChannelID, fmt.Sprintf("added song at index %d", quelen))
			if errm != nil {
				log.Println(errm)
				return
			}

		}
	case "playf":
		{
			path_to := ""
			query := strings.Join(command_args, " ")
			filepath.Walk("Y:/data/music", func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					log.Println(err)
					return nil
				}
				if info.IsDir() {
					return nil
				}
				if strings.Contains(strings.ToLower(info.Name()), strings.ToLower(query)) {
					path_to = path
					return filepath.SkipAll
				}
				return nil
			})

			if path_to == "" {
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
			entry := databaser.NewIndexEntryFromPathDnc(path_to)
			b.Queue = append(b.Queue, entry)
			b.ComChan <- bot.NewItem
			_, err := b.BState.SendMessage(c.ChannelID, fmt.Sprintf("added `%s` at index %d", entry.Title, quelen))
			if err != nil {
				log.Println(err)
				return
			}

		}
	case "playfd":
		{
			query := strings.Join(command_args, " ")
			added_count := 0
			quelen := len(b.Queue)
			filepath.WalkDir("Y:/data/music", func(path string, dir fs.DirEntry, err error) error {
				if err != nil {
					log.Println(err)
					return nil
				}
				if !dir.IsDir() {
					return nil
				}

				if strings.Contains(strings.ToLower(dir.Name()), strings.ToLower(query)) {
					filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
						if err != nil {
							log.Println(err)
						}
						if info.IsDir() {
							return nil
						}
						b.Queue = append(b.Queue, databaser.NewIndexEntryFromPathDnc(path))
						added_count++
						return nil
					})
					return filepath.SkipAll
				}
				return nil
			})

			if added_count > 0 {
				if b.VoiceSes == nil {
					utils.JoinUsersVc(b, c)
				}
				b.ComChan <- bot.NewItem
				_, err := b.BState.SendMessage(c.ChannelID, fmt.Sprintf("added `%d` items starting at index %d", added_count, quelen))
				if err != nil {
					log.Println(err)
					return
				}
			} else {
				_, err := b.BState.SendMessage(c.ChannelID, "unable to find any which match query")
				if err != nil {
					log.Println(err)
					return
				}

			}

		}
	case "pause":
		{
			if b.PlayData.Playing {
				b.PlayData.Pause()
			}
			b.PlayData.Paused = true
		}
	case "resume":
		{
			log.Println("resuming", b.PlayData.Paused)
			if b.PlayData.Paused {
				b.PlayData.Resume()
			}
		}
	case "skip":
		{
			b.PlayData.Stop()
		}
	case "restart":
		{
			b.PlayData.Restart()
		}
	case "seek":
		{
			timer, err := strconv.ParseFloat(getIndex(command_args, 0), 64)
			if err != nil {
				log.Println("failed to parse timer", err)
				return
			}
			b.PlayData.Seek(uint64(timer * 1000))
		}
	case "fuckoff":
		{
			b.PlayData.Stop()
			b.ClearQueue()
			b.VoiceSes.Leave(b.Ctx)
		}
	case "allow":
		{
			for _, mem := range command_args {
				if memid, err := strconv.Atoi(mem); err == nil {
					b.AllowList = append(b.AllowList, mem)
					go b.Db.Create(&databaser.AllowedUser{UserId: uint64(memid)})
				}
			}
		}
	case "af-set":
		{
			b.PlayData.CustomFilters = strings.Join(command_args, " ")
			log.Println("custom fitlers: ", b.PlayData.CustomFilters)
		}
	case "clear":
		{
			b.PlayData.Stop()
			b.ClearQueue()
		}
	case "remove":
		{
			idx, err := strconv.Atoi(getIndex(command_args, 0))
			if err != nil {
				b.BState.SendMessage(c.ChannelID, "what")
				return
			}
			if idx > len(b.Queue) {
				b.BState.SendMessage(c.ChannelID, "out of bounds")
				return
			}
			whatsThere := b.Queue[idx]
			b.Queue = slices.Delete(b.Queue, idx, idx+1)
			b.BState.SendMessage(c.ChannelID, fmt.Sprintf("removed `%s`", whatsThere.Title))
		}
	case "song-info":
		{
			if b.PlayData == nil {
				b.BState.SendMessage(c.ChannelID, "not playing a song baaaka")

				return
			}
			tager := b.CurrentPlayingSong()
			log.Println(tager)
			msgCnt := fmt.Sprintf("sure :3\nName: `%s`\nArtist: `%s`\nAlbum: `%s`\nLength: `%s`", tager.Title, tager.Artist, tager.Album, tager.Duration())

			b.BState.SendMessageReply(c.ChannelID, msgCnt, c.Message.ID)

		}
	case "set-qindex":
		{
			idx, err := strconv.Atoi(getIndex(command_args, 0))
			if err != nil {
				b.BState.SendMessage(c.ChannelID, "what")
				return
			}
			if idx > len(b.Queue) {
				b.BState.SendMessage(c.ChannelID, "out of bounds")
				return
			}
			if b.PlayData != nil {
				b.QueueIndex = idx - 1
				b.PlayData.Stop()
			} else {
				b.QueueIndex = idx
				b.ComChan <- bot.PlaySong
			}

		}
	case "queue":
		{
			if len(b.Queue) == 0 {
				b.BState.SendMessage(c.ChannelID, "nothing in queue")
				return
			}
			msgCnt := "queue:```"
			for k, v := range b.Queue {
				msgCnt += fmt.Sprintf("%d. %s - %s\n", k, v.Artist, v.Title)
			}
			msgCnt += "```"
			b.BState.SendMessage(c.ChannelID, msgCnt)
		}
	case "allowed":
		{
			msgCnt := "allowed:```"
			for _, v := range b.AllowList {
				msgCnt += fmt.Sprintf("%s\n", v)
			}
			msgCnt += "```"
			b.BState.SendMessage(c.ChannelID, msgCnt)
		}
	}
}

func getIndex[T any](slice []T, idx int) T {
	for k, v := range slice {
		if k == idx {
			return v
		}
	}
	var res T
	return res
}
