package commands

import (
	"fmt"
	"funnygomusic/bot"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/hako/durafmt"
	"log"
	"math"
	"time"
)

func SongInfoCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	if b.Queue.GetPlayingState() == bot.PSNotPlaying {
		b.SendMessage(c.ChannelID, "not playing a song baaaka")

		return
	}
	tager := *b.Queue.GetCurrentSong()
	log.Println(tager)
	scrobTime := b.Queue.GetCurrentSongTime().Seconds()
	msgCnt := fmt.Sprintf("sure :3\nName: `%s`\nArtist: `%s`\nAlbum: `%s`\nLength: `%s`\n%.0f%%: %.2f/%.2f",
		tager.GetTitle(), tager.GetArtist(), tager.GetAlbum(), durafmt.Parse(time.Duration(tager.GetDuration())*time.Millisecond).LimitFirstN(2),
		math.Round(scrobTime/(float64(tager.GetDuration())/1000)*100), scrobTime, float64(tager.GetDuration())/1000)

	b.SendMessageReply(c.ChannelID, msgCnt, c.Message.ID)

}
