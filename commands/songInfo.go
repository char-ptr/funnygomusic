package commands

import (
	"fmt"
	"funnygomusic/bot"
	"github.com/diamondburned/arikawa/v3/gateway"
	"log"
)

func SongInfoCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	if b.PlayData == nil {
		b.BState.SendMessage(c.ChannelID, "not playing a song baaaka")

		return
	}
	tager := b.CurrentPlayingSong()
	log.Println(tager)
	msgCnt := fmt.Sprintf("sure :3\nName: `%s`\nArtist: `%s`\nAlbum: `%s`\nLength: `%s`", tager.Title, tager.Artist, tager.Album, tager.Duration())

	b.BState.SendMessageReply(c.ChannelID, msgCnt, c.Message.ID)

}
