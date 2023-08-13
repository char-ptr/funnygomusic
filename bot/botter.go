package bot

import (
	"context"
	"funnygomusic/databaser"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/voice"
	"gorm.io/gorm"
)

type ComData int

type Botter struct {
	BState    *state.State
	VoiceSes  *voice.Session
	Queue     *QueueManager
	Ctx       context.Context
	MyId      discord.UserID
	AllowList []string
	SubChan   discord.ChannelID
	Db        *gorm.DB
}

func NewBotter(s *state.State, ctx *context.Context) *Botter {
	b := &Botter{
		BState:    s,
		Ctx:       *ctx,
		AllowList: []string{},
		Db:        databaser.NewDatabase(),
	}
	b.Queue = NewQueueManager(b)
	return b

}
