package commands

import (
	"fmt"
	"funnygomusic/bot"
	"github.com/diamondburned/arikawa/v3/gateway"
	"log"
	"math"
)

func SongInfoCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	if b.Queue.GetPlayingState() == bot.PSNotPlaying {
		b.SendMessage(c.ChannelID, "not playing a song baaaka")

		return
	}
	tager := b.Queue.GetCurrentSong()
	log.Println(tager)
	scrobTime := b.Queue.GetCurrentSongTime().Seconds()
	msgCnt := fmt.Sprintf("sure :3\nName: `%s`\nArtist: `%s`\nAlbum: `%s`\nLength: `%s`\n%.0f%%: %.2f/%.2f", tager.Title, tager.Artist, tager.Album, tager.DurationStr(), math.Round(scrobTime/tager.Length*100), scrobTime, tager.Length)

	b.SendMessageReply(c.ChannelID, msgCnt, c.Message.ID)

}
