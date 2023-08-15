package commands

import (
	"funnygomusic/bot"
	"github.com/diamondburned/arikawa/v3/gateway"
	"log"
)

func TestVceCommand(c *gateway.MessageCreateEvent, b *bot.Botter, args []string) {
	log.Printf("TestVceCommand: %#v", b.VoiceSes)
}
