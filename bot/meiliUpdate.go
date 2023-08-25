package bot

import (
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

func (b *Botter) MeiliUpdate() {

	sidx := b.Meili.Index("songs")
	//aridx := b.Meili.Index("artists")
	//alidx := b.Meili.Index("albums")
	var all_songs []IndexedSongExt
	b.Db.Raw("select iso.*, iaa.id as artistID, ia.id as albumID, iaa.name as artist, ia.name as album from indexed_songs iso join public.indexed_albums ia on iso.album = ia.id join indexed_artists iaa on iso.artist = iaa.id where iso.id != ''").Scan(&all_songs)

	//for _, song := range all_songs {
	//	//log.Println(song.Title, song.Artist, song.Album)
	//}
	log.Println("songs", len(all_songs))

	tinfo, err := sidx.AddDocuments(all_songs)
	if err != nil {
		log.Println("failed to add songs to meili", err)
	}
	log.Println("meili request ok: ", tinfo.Status, tinfo.Type)
}
