package bot

import (
	"context"
	"log"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/voice"
	"github.com/diamondburned/arikawa/v3/voice/udp"
	"github.com/diamondburned/arikawa/v3/voice/voicegateway"
	"github.com/pkg/errors"
)

type VoiceSessionUser struct {
	Name     string         `json:"name"`
	ID       discord.UserID `json:"id"`
	Muted    bool           `json:"muted"`
	Deafen   bool           `json:"deafen"`
	Speaking bool           `json:"speaking"`
}
type VoiceSessionHndlr struct {
	vs        *voice.Session
	Users     []VoiceSessionUser `json:"users"`
	GuildId   discord.GuildID    `json:"guild_id"`
	ChannelID discord.ChannelID  `json:"channel_id"`
}

func NewVoiceSessionUser(vste *discord.VoiceState, b *Botter) VoiceSessionUser {
	uname := b.MyUsername
	if vste.UserID != b.MyId {
		u, err := b.User(vste.UserID)
		if err != nil {
			uname = "unknown"
		} else {
			uname = u.DisplayOrUsername()
		}
	}
	return VoiceSessionUser{
		Name:     uname,
		ID:       vste.UserID,
		Muted:    vste.Mute || vste.SelfMute,
		Deafen:   vste.Deaf || vste.SelfDeaf,
		Speaking: false,
	}
}

const (
	FrameDuration = 60 // ms
	TimeIncrement = 2880
)

func (v *VoiceSessionHndlr) AttachVoiceSession(vs *voice.Session) {
	v.vs = vs
	//vs.AddHandler(func(spk *voicegateway.SpeakingEvent) {
	//	log.Println("spk evt", spk)
	//	for k, user := range v.Users {
	//		if user.ID == spk.UserID {
	//			log.Println("spk = ", spk.Speaking, voicegateway.Microphone)
	//			v.Users[k].Speaking = spk.Speaking == voicegateway.Microphone
	//		}
	//	}
	//})
}

func (v *VoiceSessionHndlr) JoinedChannel(chn *discord.ChannelID, gld *discord.GuildID) {
	v.GuildId = *gld
	v.ChannelID = *chn
}

func (v *VoiceSessionHndlr) UpdateUsers(b *Botter) {
	v.Users = nil
	vstates, err := b.VoiceStates(v.GuildId)
	if err != nil {
		log.Println("unable to getg voice states: ", err)
		return
	}
	for _, state := range vstates {
		if state.ChannelID != v.ChannelID {
			continue
		}
		v.Users = append(v.Users, NewVoiceSessionUser(&state, b))
	}
}

func (v *VoiceSessionHndlr) JoinUsersVc(b *Botter, gld discord.GuildID, uid discord.UserID) error {
	vs, err := voice.NewSession(b)
	if err != nil {
		return errors.Wrap(err, "cannot make new voice session")
	}
	v.AttachVoiceSession(vs)
	vs.SetUDPDialer(udp.DialFuncWithFrequency(
		FrameDuration*time.Millisecond, // correspond to -frame_duration
		TimeIncrement,
	))
	uservs, err := b.VoiceState(gld, uid)
	v.GuildId = gld
	v.ChannelID = uservs.ChannelID

	if err != nil {
		return errors.Wrap(err, "cannot get voice state")
	}
	vs.JoinChannel(b.Ctx, uservs.ChannelID, false, false)
	go v.UpdateUsers(b)
	return nil
}

func (v *VoiceSessionHndlr) Open() bool {
	return v.vs != nil
}

func (v *VoiceSessionHndlr) Leave(ctx context.Context) {
	v.vs.Leave(ctx)
	v.ChannelID = discord.ChannelID(0)
	v.GuildId = discord.GuildID(0)
	v.Users = nil
	v.vs = nil
}

func (v *VoiceSessionHndlr) HasUser(uid discord.UserID) bool {
	for _, user := range v.Users {
		if user.ID == uid {
			return true
		}
	}
	return false
}

func (v *VoiceSessionHndlr) DeleteUser(uid discord.UserID) {
	for k, user := range v.Users {
		if user.ID == uid {
			v.Users = append(v.Users[:k], v.Users[k+1:]...)
			return
		}
	}
}

func (v *VoiceSessionHndlr) GetUser(uid discord.UserID) *VoiceSessionUser {
	for _, user := range v.Users {
		if user.ID == uid {
			return &user
		}
	}
	return nil
}

func (v *VoiceSessionHndlr) UpdateUser(uid discord.UserID, user *VoiceSessionUser) {
	for k, u := range v.Users {
		if u.ID == uid {
			v.Users[k] = *user
			return
		}
	}
}

func (v *VoiceSessionHndlr) AddUser(user *VoiceSessionUser) {
	v.Users = append(v.Users, *user)
}

func (v *VoiceSessionHndlr) Speaking(isSpeaking bool, ctx context.Context) {
	if !v.Open() {
		return
	}
	if isSpeaking {
		v.vs.Speaking(ctx, voicegateway.Microphone)
	} else {
		v.vs.Speaking(ctx, voicegateway.NotSpeaking)
	}
}

func (v *VoiceSessionHndlr) GetSession() *voice.Session {
	return v.vs
}
