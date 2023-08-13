package routes

import (
	"funnygomusic/bot"
	"github.com/gin-gonic/gin"
)

func Connected(c *gin.Context) {
	b := c.MustGet("bot").(*bot.Botter)
	c.JSON(200, gin.H{
		"Connected": b.V.Open(),
	})
}
