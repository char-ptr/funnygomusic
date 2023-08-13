package bot

import (
	"context"
	"funnygomusic/databaser"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"gorm.io/gorm"
)

type ComData int

type Botter struct {
	BState     *state.State
	V          *VoiceSessionHndlr
	Queue      *QueueManager
	Ctx        context.Context
	MyId       discord.UserID
	MyUsername string
	AllowList  []string
	SubChan    discord.ChannelID
	Db         *gorm.DB
}

func NewBotter(s *state.State, ctx *context.Context) *Botter {
	b := &Botter{
		BState:    s,
		Ctx:       *ctx,
		AllowList: []string{},
		V:         &VoiceSessionHndlr{},
		Db:        databaser.NewDatabase(),
	}
	b.V.b = b
	b.Queue = NewQueueManager(b)
	return b

}
