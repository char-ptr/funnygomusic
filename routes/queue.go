package routes

import (
	"database/sql"
	_ "embed"
	"funnygomusic/bot"
	"funnygomusic/bot/entries"
	"github.com/gin-gonic/gin"
	"log"
)

func CurrentSong(c *gin.Context) {
	b := c.MustGet("bot").(*bot.Botter)
	if b.VoiceSes.Open() == false {
		c.JSON(200, nil)
	}
	c.JSON(200, bot.GetTypedEntry(b.Queue.GetCurrentSong()))
}
func GetQueue(c *gin.Context) {
	b := c.MustGet("bot").(*bot.Botter)
	if b.VoiceSes.Open() == false {
		c.JSON(200, nil)
	}
	c.JSON(200, b.Queue.GetEntries())
}
func PushToQueue(c *gin.Context) {
	b := c.MustGet("bot").(*bot.Botter)
	if b.VoiceSes.Open() == false {
		c.JSON(200, nil)
		return
	}
	var pe []PushEntry
	err := c.BindJSON(&pe)
	if err != nil {
		c.JSON(400, err)
		return
	}
	var peIds []string
	for _, v := range pe {
		peIds = append(peIds, v.Id)
	}
	/// what the fuck?! - basically just selects basic song fields from either a: (song, album, or artist) id.
	rows, err := b.Db.Raw(songFetcherQuery,
		sql.Named("ids", peIds),
	).Rows()
	defer rows.Close()
	for rows.Next() {
		var songer entries.Indexed
		err := b.Db.ScanRows(rows, &songer)
		if err != nil {
			c.JSON(500, err)
			return
		}
		log.Println("request to add", songer.GetTitle())
		b.Queue.Notify <- bot.NewPlaylistMessage(bot.PlaylistAdd).SetEntry(&songer)
	}

	c.JSON(200, gin.H{"ok": "ok"})
}

//go:embed songFetcher.sql
var songFetcherQuery string

type PushEntry struct {
	Id string
}
