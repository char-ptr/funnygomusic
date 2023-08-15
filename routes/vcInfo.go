package routes

import (
	"funnygomusic/bot"
	"github.com/gin-gonic/gin"
)

func VcInfo(c *gin.Context) {
	b := c.MustGet("bot").(*bot.Botter)
	c.JSON(200, b.VoiceSes)
}
