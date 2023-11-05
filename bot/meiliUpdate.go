package bot

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"log"
)

type IndexedSongExt struct {
	gorm.Model
	Title    string
	Artist   string
	Album    string
	AlbumID  string
	ArtistID string
	Track    int
	File     string
	Length   float64
	ID       string `gorm:"primaryKey"`
}
type AlbumExt struct {
	gorm.Model
	Name     string
	Artist   string
	ArtistID string
	Tracks   pq.StringArray `gorm:"type:text[]"`
	ID       string         `gorm:"primaryKey"`
}

func (b *Botter) MeiliUpdate() {

	sidx := b.Meili.Index("songs")
	//aridx := b.Meili.Index("artists")
	alidx := b.Meili.Index("albums")
	var all_songs []IndexedSongExt
	var all_albums []AlbumExt
	b.Db.Raw("select iso.*, iso.artist as artistID, iso.album as albumID, iaa.name as artist, ia.name as album from indexed_songs iso join public.indexed_albums ia on iso.album = ia.id join indexed_artists iaa on iso.artist = iaa.id where iso.id != ''").Scan(&all_songs)
	b.Db.Raw("select ia.name, iaa.name as artist, ia.artist as artistID,ia.id, array(select iso.title from indexed_songs iso where iso.album = ia.id) as tracks from indexed_albums ia join indexed_artists iaa on iaa.id = ia.artist;").Scan(&all_albums)
	if all_songs == nil {
		log.Println("no songs found")
		return
	}
	tinfo, err := sidx.AddDocuments(all_songs, "ID")
	alidx.AddDocuments(all_albums, "ID")
	if err != nil {
		log.Println("failed to add songs to meili", err)
	}
	log.Println("meili request ok: ", tinfo.Status, tinfo.Type)
}
