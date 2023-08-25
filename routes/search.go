package routes

import (
	"funnygomusic/bot"
	"github.com/gin-gonic/gin"
)

func SearchForSong(c *gin.Context) {
	b := c.MustGet("bot").(*bot.Botter)

	b.Queue.Notify <- bot.NewPlaylistMessage(bot.CurrentStop)
	c.JSON(200, gin.H{"ok": "ok"})
}
