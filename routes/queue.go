package routes

import (
	"database/sql"
	_ "embed"
	"funnygomusic/bot"
	"funnygomusic/bot/entries"
	"log"

	"github.com/gin-gonic/gin"
)

func CurrentSong(c *gin.Context) {
	b := c.MustGet("bot").(*bot.Botter)
	if !b.VoiceSes.Open() {
		c.JSON(200, nil)
	}
	c.JSON(200, bot.GetTypedEntry(b.Queue.GetCurrentSong()))
}

func GetQueue(c *gin.Context) {
	b := c.MustGet("bot").(*bot.Botter)
	if !b.VoiceSes.Open() {
		c.JSON(200, nil)
	}
	c.JSON(200, gin.H{"queue": b.Queue.GetEntries(), "index": b.Queue.GetIndex(), "length": len(b.Queue.GetEntries())})
}

func ClearQueue(c *gin.Context) {
	b := c.MustGet("bot").(*bot.Botter)
	if !b.VoiceSes.Open() {
		c.JSON(200, nil)
	}
	b.Queue.Notify <- bot.NewPlaylistMessage(bot.PlaylistClear)
	c.JSON(200, gin.H{"ok": "ok"})
}

func SkipSong(c *gin.Context) {
	b := c.MustGet("bot").(*bot.Botter)
	if !b.VoiceSes.Open() {
		c.JSON(200, nil)
	}
	b.Queue.Notify <- bot.NewPlaylistMessage(bot.CurrentSkip)
	c.JSON(200, gin.H{"ok": "ok"})
}

func PauseSong(c *gin.Context) {
	b := c.MustGet("bot").(*bot.Botter)
	if !b.VoiceSes.Open() {
		c.JSON(200, nil)
	}
	b.Queue.Notify <- bot.NewPlaylistMessage(bot.CurrentPause)
	c.JSON(200, gin.H{"ok": "ok"})
}

func ResumeSong(c *gin.Context) {
	b := c.MustGet("bot").(*bot.Botter)
	if !b.VoiceSes.Open() {
		c.JSON(200, nil)
	}
	b.Queue.Notify <- bot.NewPlaylistMessage(bot.CurrentResume)
	c.JSON(200, gin.H{"ok": "ok"})
}

func StopSong(c *gin.Context) {
	b := c.MustGet("bot").(*bot.Botter)
	if !b.VoiceSes.Open() {
		c.JSON(200, nil)
	}
	b.Queue.Notify <- bot.NewPlaylistMessage(bot.CurrentStop)
	c.JSON(200, gin.H{"ok": "ok"})
}

func SeekSong(c *gin.Context) {
	b := c.MustGet("bot").(*bot.Botter)
	if !b.VoiceSes.Open() {
		c.JSON(200, nil)
	}
	var seek struct {
		Position int `json:"position"`
	}
	err := c.BindJSON(&seek)
	if err != nil {
		c.JSON(400, err)
		return
	}
	b.Queue.Notify <- bot.NewPlaylistMessage(bot.CurrentSkip).SetSeek(seek.Position)
	c.JSON(200, gin.H{"ok": "ok"})
}

func PushToQueue(c *gin.Context) {
	b := c.MustGet("bot").(*bot.Botter)
	if !b.VoiceSes.Open() {
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
