package commands

import (
	"fmt"
	"funnygomusic/bot"
	"github.com/diamondburned/arikawa/v3/gateway"
	"log"
)

func SongInfoCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	if b.Queue.GetPlayingState() == bot.PSNotPlaying {
		b.BState.SendMessage(c.ChannelID, "not playing a song baaaka")

		return
	}
	tager := b.Queue.GetCurrentSong()
	log.Println(tager)
	msgCnt := fmt.Sprintf("sure :3\nName: `%s`\nArtist: `%s`\nAlbum: `%s`\nLength: `%s`\n%.2f/%.2f", tager.Title, tager.Artist, tager.Album, tager.Duration(), b.Queue.GetCurrentSongTime().Seconds(), tager.Length.Seconds())

	b.BState.SendMessageReply(c.ChannelID, msgCnt, c.Message.ID)

}
