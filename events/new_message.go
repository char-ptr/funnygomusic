package events

import (
	"fmt"
	"funnygomusic/bot"
	"funnygomusic/utils"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dhowden/tag"
	"github.com/diamondburned/arikawa/v3/gateway"
	"golang.org/x/exp/slices"
)

func onMessage(c *gateway.MessageCreateEvent, leBot *bot.Botter) {
	if !strings.HasPrefix(c.Message.Content, "`") || !(slices.Contains(leBot.AllowList, c.Author.ID.String())) {
		return
	}
	command_args := strings.Split(c.Message.Content, " ")
	command := command_args[0][1:]
	command_args = command_args[1:]
	leBot.SubChan = c.ChannelID
	switch command {
	case "join":
		{
			utils.JoinUsersVc(leBot, c)
		}

	case "play":
		{
			// if leBot.PlayData != nil && leBot.PlayData.Playing {
			// 	log.Println("already playing")
			// 	return
			// }
			path_to := strings.Join(command_args, " ")
			fInfo, err := os.Stat(path_to)
			if os.IsNotExist(err) {
				log.Println("file does not exist")
				return
			}
			if leBot.VoiceSes == nil {
				utils.JoinUsersVc(leBot, c)
			}
			quelen := len(leBot.Queue)
			leBot.Queue = append(leBot.Queue, bot.QueueEntry{Path: path_to, Name: fInfo.Name()})
			leBot.ComChan <- bot.NewItem
			_, errm := leBot.BState.SendMessage(c.ChannelID, fmt.Sprintf("added song at index %d", quelen))
			if err != nil {
				log.Println(errm)
				return
			}

		}
	case "playf":
		{
			path_to := ""
			var fInfo fs.FileInfo = nil
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
					fInfo = info
					return filepath.SkipAll
				}
				return nil
			})

			if path_to == "" {
				_, err := leBot.BState.SendMessage(c.ChannelID, "unable to find, srry")
				if err != nil {
					log.Println(err)
					return
				}
				return

			}
			if leBot.VoiceSes == nil {
				utils.JoinUsersVc(leBot, c)
			}
			quelen := len(leBot.Queue)
			leBot.Queue = append(leBot.Queue, bot.QueueEntry{Path: path_to, Name: fInfo.Name()})
			leBot.ComChan <- bot.NewItem
			_, err := leBot.BState.SendMessage(c.ChannelID, fmt.Sprintf("added `%s` at index %d", fInfo.Name(), quelen))
			if err != nil {
				log.Println(err)
				return
			}

		}
	case "pause":
		{
			if leBot.PlayData.Playing {
				leBot.PlayData.Pause()
			}
			leBot.PlayData.Paused = true
		}
	case "resume":
		{
			log.Println("resuming", leBot.PlayData.Paused)
			if leBot.PlayData.Paused {
				leBot.PlayData.Resume()
			}
		}
	case "skip":
		{
			leBot.PlayData.Stop()
		}
	case "restart":
		{
			leBot.PlayData.Restart()
		}
	case "seek":
		{
			timer, err := strconv.ParseFloat(getIndex(command_args, 0), 64)
			if err != nil {
				log.Println("failed to parse timer", err)
				return
			}
			leBot.PlayData.Seek(uint64(timer * 1000))
		}
	case "fuckoff":
		{
			leBot.PlayData.Stop()
			leBot.ClearQueue()
			leBot.VoiceSes.Leave(leBot.Ctx)
		}
	case "allow":
		{
			for _, mem := range c.Mentions {
				leBot.AllowList = append(leBot.AllowList, mem.ID.String())
			}
		}
	case "af-set":
		{
			leBot.PlayData.CustomFilters = strings.Join(command_args, " ")
			log.Println("custom fitlers: ", leBot.PlayData.CustomFilters)
		}
	case "clear":
		{
			leBot.PlayData.Stop()
			leBot.ClearQueue()
		}
	case "remove":
		{
			idx, err := strconv.Atoi(getIndex(command_args, 0))
			if err != nil {
				leBot.BState.SendMessage(c.ChannelID, "what")
				return
			}
			if idx > len(leBot.Queue) {
				leBot.BState.SendMessage(c.ChannelID, "out of bounds")
				return
			}
			whatsThere := leBot.Queue[idx]
			leBot.Queue = slices.Delete(leBot.Queue, idx, idx+1)
			leBot.BState.SendMessage(c.ChannelID, fmt.Sprintf("removed `%s`", whatsThere.Name))
		}
	case "song-info":
		{
			if leBot.PlayData == nil {
				leBot.BState.SendMessage(c.ChannelID, "not playing a song baaaka")

				return
			}
			thef, _ := os.Open(leBot.CurrentPlayingSong().Path)
			defer thef.Close()
			tag, err := tag.ReadFrom(thef)
			if err != nil {
				leBot.BState.SendMessage(c.ChannelID, "error opening file")
				return
			}
			msg_cnt := fmt.Sprintf("sure :3\nName: `%s`\nArtist: `%s`\nAlbum: `%s`", tag.Title(), tag.Artist(), tag.Album())

			leBot.BState.SendMessageReply(c.ChannelID, msg_cnt, c.Message.ID)

		}
	}
}
func OnMessage(le_bot *bot.Botter) func(c *gateway.MessageCreateEvent) {
	return func(c *gateway.MessageCreateEvent) {
		onMessage(c, le_bot)
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
