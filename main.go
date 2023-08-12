package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"funnygomusic/bot"
	"funnygomusic/events"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

var Bot bot.Botter = bot.Botter{}

func main() {

	id := gateway.DefaultIdentifier(os.Getenv("BOT_TOKEN"))
	id.Properties.OS = "iOS"
	id.Properties.Browser = "Discord iOS"

	current_state := state.NewWithIdentifier(id)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	Bot = bot.NewBotter(current_state, &ctx)
	Bot.AllowList = append(Bot.AllowList, os.Getenv("BOT_OWNER"))
	go Bot.PlayManagerStart()

	current_state.AddHandler(events.OnMessage(&Bot))
	current_state.AddHandler(events.OnReady(&Bot))

	if err := current_state.Open(ctx); err != nil {
		log.Fatalln("failed to open:", err)
	}
	defer current_state.Close()

	log.Println("Blocking, press ctrl+c to continue...")
	<-ctx.Done()
}
