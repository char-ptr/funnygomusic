package bot

import (
	"context"
	"funnygomusic/databaser"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/meilisearch/meilisearch-go"
	"gorm.io/gorm"
)

type ComData int

type Botter struct {
	Ctx context.Context
	*state.State
	VoiceSes   *VoiceSessionHndlr
	Queue      *QueueManager
	Db         *gorm.DB
	Meili      *meilisearch.Client
	MyUsername string
	AllowList  []string
	MyId       discord.UserID
	SubChan    discord.ChannelID
}

func NewBotter(s *state.State, ctx *context.Context) *Botter {
	b := &Botter{
		State:     s,
		Ctx:       *ctx,
		AllowList: []string{},
		VoiceSes:  &VoiceSessionHndlr{},
		Db:        databaser.NewDatabase(),
		Meili:     databaser.NewMeili(),
	}
	b.MeiliUpdate()
	b.Queue = NewQueueManager(b)
	return b
}
