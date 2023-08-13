package commands

import (
	"funnygomusic/bot"
	"github.com/diamondburned/arikawa/v3/gateway"
)

type Command = func(*gateway.MessageCreateEvent, *bot.Botter, []string)
