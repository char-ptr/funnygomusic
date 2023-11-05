package entries

import (
	"funnygomusic/bot"
	"funnygomusic/bot/players"
	"gorm.io/gorm"
)

type Indexed struct {
	F        string
	Title    string
	Artist   string
	Album    string
	Duration float64
	ID       string
}

func (l *Indexed) GetTitle() string {
	return l.Title
}
func (l *Indexed) GetAlbum() string {
	return l.Album
}
func (l *Indexed) GetArtist() string {
	return l.Artist
}
func (l *Indexed) GetDuration() int {
	return int(l.Duration * 1000)
}
func (l *Indexed) GetPlayer() bot.Player {
	return players.NewFilePlayer(l.F)
}
func (l *Indexed) GetID() string {
	return l.ID
}
func NewIndexedFromDb(trkID string, db *gorm.DB) (idx Indexed) {
	db.Raw("select file as f, title, il.name as album, ia.name as artist, iso.id as id, iso.length as duration from indexed_songs as iso left join indexed_albums as il on il.id = iso.album left join indexed_artists as ia on ia.id = iso.artist where iso.id = ?", trkID).Scan(&idx)
	return
}
