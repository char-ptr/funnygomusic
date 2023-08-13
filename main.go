package main

import (
	"context"
	"funnygomusic/databaser"
	"log"
	"os"
	"os/signal"
	"strconv"

	"funnygomusic/bot"
	_ "funnygomusic/databaser"
	"funnygomusic/events"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

func main() {

	id := gateway.DefaultIdentifier(os.Getenv("BOT_TOKEN"))
	id.Properties.OS = "iOS"
	id.Properties.Browser = "Discord iOS"
	currentState := state.NewWithIdentifier(id)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	b := bot.NewBotter(currentState, &ctx)
	var alusers []databaser.AllowedUser
	b.Db.Find(&alusers)
	b.AllowList = append(b.AllowList, os.Getenv("BOT_OWNER"))
	for _, u := range alusers {
		b.AllowList = append(b.AllowList, strconv.Itoa(int(u.UserId)))
	}

	go b.PlayManagerStart()

	currentState.AddHandler(func(c *gateway.ReadyEvent) { events.OnReady(c, &b) })
	currentState.AddHandler(func(c *gateway.MessageCreateEvent) { events.OnMessage(c, &b) })

	// start connection
	if err := currentState.Open(ctx); err != nil {
		log.Fatalln("failed to open:", err)
	}
	defer currentState.Close()

	log.Println("Blocking, press ctrl+c to continue...")
	<-ctx.Done()
}
